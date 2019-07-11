package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v2"
	"strconv"
	"time"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/People/NewStruct/g
  :%s/people/newinst/g

**/

// XXX custom struct name and fields
type People struct {
	Id           int64     `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	XAddress     string    `db:"xaddress,size:65534" json:"xaddress"`
	XDateOfBirth string    `db:"xdob" json:"xdob"`
	Role         string    `db:"role" json:"role"`
	Status       string    `db:"status" json:"status"`
	Created      time.Time `db:"created" json:"created"` // or int64
	Updated      time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func (a *People) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now() // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

func (a *People) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	return nil
}

// REST handlers

func GetPeoples(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	query := "SELECT * FROM people"

	// Parse query string
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	var count int64
	if s != "" {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM people  WHERE " + s)
		query = query + " WHERE " + s
	} else {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM people")
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

	var peoples []People
	_, err := dbmap.Select(&peoples, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10)) // float64 to string
		c.JSON(200, peoples)
	} else {
		c.JSON(404, gin.H{"error": "no people(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/peoples
}

func GetPeople(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var people People
	err := dbmap.SelectOne(&people, "SELECT * FROM people WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, people)
	} else {
		c.JSON(404, gin.H{"error": "people not found"})
	}

	// curl -i http://localhost:8080/api/v1/peoples/1
}

func PostPeople(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)

	var people People
	c.Bind(&people)

	if verbose == true {
		fmt.Println(people)
	}

	// XXX Check encrypted field match a valid key
	if people.XAddress != "" && TestValidSalt(dbmap, people.XAddress[0:16]) == false {
		people.Name = ""
	}
	if people.XDateOfBirth != "" && TestValidSalt(dbmap, people.XDateOfBirth[0:16]) == false {
		people.Name = ""
	}

	if people.Name != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&people)
		if err == nil {
			c.JSON(201, people)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/peoples
}

func UpdatePeople(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	id := c.Params.ByName("id")

	var people People
	err := dbmap.SelectOne(&people, "SELECT * FROM people WHERE id=?", id)
	if err == nil {
		var json People
		c.Bind(&json)
		if verbose == true {
			fmt.Println(json)
		}
		people_id, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		people := People{
			Id:           people_id,
			Name:         json.Name,
			XAddress:     json.XAddress,
			XDateOfBirth: json.XDateOfBirth,
			Role:         json.Role,
			Status:       json.Status,
			Created:      people.Created, //people read from previous select
		}

		// XXX Check encrypted field match a valid key
		// else create mandatory field error
		if people.XAddress != "" && TestValidSalt(dbmap, people.XAddress[0:16]) == false {
			people.Name = ""
		}
		if people.XDateOfBirth != "" && TestValidSalt(dbmap, people.XDateOfBirth[0:16]) == false {
			people.Name = ""
		}

		if people.Name != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&people)
			if err == nil {
				c.JSON(200, people)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "people not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/peoples/1
}

func DeletePeople(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var people People
	err := dbmap.SelectOne(&people, "SELECT * FROM people WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&people)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "people not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/peoples/1
}
