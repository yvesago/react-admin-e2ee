package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v2"
	//	"strconv"
	"time"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/Util/NewStruct/g
  :%s/util/newinst/g

**/

// XXX custom struct name and fields
type Util struct {
	Id        int64     `db:"id" json:"id"`
	VerifyKey string    `db:"verifykey" json:"verifykey"`
	Created   time.Time `db:"created" json:"created"` // or int64
	Updated   time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func (a *Util) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now() // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

func (a *Util) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	CurrentSalt = a.VerifyKey[0:16]
	return nil
}

// REST handlers

func GetVerifKey(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	//id := c.Params.ByName("id")

	var util Util
	err := dbmap.SelectOne(&util, "SELECT * FROM util WHERE id=1 LIMIT 1")

	if err == nil {
		c.JSON(200, util)
	} else {
		c.JSON(404, gin.H{"error": "util not found"})
	}

	// curl -i http://localhost:8080/api/v1/utils/1
}

func UpdateVerifKey(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	//id := c.Params.ByName("id")

	var util Util
	err := dbmap.SelectOne(&util, "SELECT * FROM util WHERE id=1")
	if err == nil {
		var json Util
		c.Bind(&json)

		if verbose == true {
			fmt.Println(json)
		}

		//util_id, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		util := Util{
			Id:        1,
			VerifyKey: json.VerifyKey,
			Created:   util.Created, //util read from previous select
		}

		if util.Id != 0 { // XXX Check mandatory fields
			_, err = dbmap.Update(&util)
			if err == nil {
				c.JSON(200, util)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "util not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/utils/1
}

/*
func GetUtils(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	query := "SELECT * FROM util"

	// Parse query string
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	var count int64
	if s != "" {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM util  WHERE " + s)
		query = query + " WHERE " + s
	} else {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM util")
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	if verbose == true {
		fmt.Println(q)
		fmt.Println(" -- " + query)
	}

	var utils []Util
	_, err := dbmap.Select(&utils, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10)) // float64 to string
		c.JSON(200, utils)
	} else {
		c.JSON(404, gin.H{"error": "no util(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/utils
}

func GetUtil(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var util Util
	err := dbmap.SelectOne(&util, "SELECT * FROM util WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, util)
	} else {
		c.JSON(404, gin.H{"error": "util not found"})
	}

	// curl -i http://localhost:8080/api/v1/utils/1
}

func PostUtil(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)

	var util Util
	c.Bind(&util)

	if verbose == true {
		fmt.Println(util)
	}

	if util.Id != 0 { // XXX Check mandatory fields
		err := dbmap.Insert(&util)
		if err == nil {
			c.JSON(201, util)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/utils
}

func UpdateUtil(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	id := c.Params.ByName("id")

	var util Util
	err := dbmap.SelectOne(&util, "SELECT * FROM util WHERE id=?", id)
	if err == nil {
		var json Util
		c.Bind(&json)

		if verbose == true {
			fmt.Println(json)
		}

		util_id, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		util := Util{
			Id:        util_id,
			VerifyKey: json.VerifyKey,
			Created:   util.Created, //util read from previous select
		}

		if util.Id != 0 { // XXX Check mandatory fields
			_, err = dbmap.Update(&util)
			if err == nil {
				c.JSON(200, util)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "util not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/utils/1
}

func DeleteUtil(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var util Util
	err := dbmap.SelectOne(&util, "SELECT * FROM util WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&util)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "util not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/utils/1
}
*/
