package models

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// TestValidSalt : test if key is a valid salt in vaults
func TestValidSalt(dbmap *gorp.DbMap, test string) bool {
        var vault Vault
        err := dbmap.SelectOne(&vault, "SELECT id FROM vault WHERE verifykey LIKE ?", test+"%")
        if err == nil && vault.Id != 0 {
            return true
        }
        return false
}

// Database : gin Middlware to select database
func Database(connString string) gin.HandlerFunc {
	dbmap := InitDb(connString)
	return func(c *gin.Context) {
		c.Set("DBmap", dbmap)
		c.Next()
	}
}

// InitDb : create or update and connect to db on startup
func InitDb(dbName string) *gorp.DbMap {
	// XXX fix database type
	db, err := sql.Open("sqlite3", dbName)
	checkErr(err, "sql.Open failed")
	//dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	// XXX fix tables names
	dbmap.AddTableWithName(People{}, "People").SetKeys(true, "Id")
	dbmap.AddTableWithName(User{}, "User").SetKeys(true, "Id")
	dbmap.AddTableWithName(Vault{}, "Vault").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	/*var u Util
	dbmap.SelectOne(&u, "select * from Util where id = 1")
	if u.Id != 1 {
		dbmap.Insert(&u)
	} else {
		CurrentSalt = u.VerifyKey[0:16]
	}*/

	return dbmap
}


// ParseQuery : Parse a http query
func ParseQuery(q map[string][]string) (string, string, string) {
	query := ""
	if q["_filters"] != nil {
		data := make(map[string]string)
		err := json.Unmarshal([]byte(q["_filters"][0]), &data)
		if err == nil {
			var searches []string
			for col, search := range data {
				valid := regexp.MustCompile("^[A-Za-z0-9_.]+$")
				if col != "" && search != "" && valid.MatchString(col) && valid.MatchString(search) {
					searches = append(searches, col+" LIKE \"%"+search+"%\"")
				}
			}
			query = query + strings.Join(searches, " AND ") // TODO join with OR for same keys
		}
	}

	sort := ""
	if q["_sortField"] != nil && q["_sortDir"] != nil {
		sortField := q["_sortField"][0]
		// prevent SQLi
		valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
		if !valid.MatchString(sortField) {
			sortField = ""
		}
		if sortField == "created" || sortField == "updated" { // XXX trick for sqlite
			sortField = "datetime(" + sortField + ")"
		}
		sortOrder := q["_sortDir"][0]
		if sortOrder != "ASC" {
			sortOrder = "DESC"
		}
		if sortField != "" {
			sort = " ORDER BY " + sortField + " " + sortOrder
		}
	}

	limit := ""
	// _page, _perPage : LIMIT + OFFSET
	perPageInt := 0
	if q["_perPage"] != nil {
		perPage := q["_perPage"][0]
		valid := regexp.MustCompile("^[0-9]+$")
		if valid.MatchString(perPage) {
			perPageInt, _ = strconv.Atoi(perPage)
			limit = " LIMIT " + perPage
		}
	}
	if q["_page"] != nil {
		page := q["_page"][0]
		valid := regexp.MustCompile("^[0-9]+$")
		pageInt, _ := strconv.Atoi(page)

		if valid.MatchString(page) && pageInt > 1 {
			offset := (pageInt-1)*perPageInt + 1
			limit = limit + " OFFSET " + strconv.Itoa(offset)
		}
	}

	// _start, _end : LIMIT start, size
	if q["_start"] != nil && q["_end"] != nil {
		start := q["_start"][0]
		end := q["_end"][0]
		valid := regexp.MustCompile("^[0-9]+$")
		startInt, _ := strconv.Atoi(start)
		endInt, _ := strconv.Atoi(end)
		startInt = startInt - 1 // indice start from 0

		if valid.MatchString(start) && valid.MatchString(end) && endInt > startInt {
			size := endInt - startInt
			limit = " LIMIT " + strconv.Itoa(startInt) + ", " + strconv.Itoa(size)
		}
	}

	return query, sort, limit
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
