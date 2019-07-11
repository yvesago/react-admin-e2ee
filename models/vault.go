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
  :%s/Vault/NewStruct/g
  :%s/vault/newinst/g

**/

// Vault fields XXX custom struct name and fields
type Vault struct {
	Id        int64  `db:"id" json:"id"`
	VaultName string `db:"vaultname" json:"vaultname"`
	VerifyKey string `db:"verifykey" json:"verifykey"`
	// TODO:
	//  User
	//  UserEnckey
	//  Role  (su, admin, user, reader)
	//  Fragment (for su sss)
	Created time.Time `db:"created" json:"created"` // or int64
	Updated time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

// PreInsert : set create and update time
func (a *Vault) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now() // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

// PreUpdate : set update time
func (a *Vault) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	//CurrentSalt = a.VerifyKey[0:16]
	return nil
}

// REST handlers

// GetVaults : get all vaults
func GetVaults(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	query := "SELECT * FROM vault"

	// Parse query string
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	var count int64
	if s != "" {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM vault  WHERE " + s)
		query = query + " WHERE " + s
	} else {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM vault")
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

	var vaults []Vault
	_, err := dbmap.Select(&vaults, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10)) // float64 to string
		c.JSON(200, vaults)
	} else {
		c.JSON(404, gin.H{"error": "no vault(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/vaults
}

// GetVault : mainly to get VerifyKey
func GetVault(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var vault Vault
	err := dbmap.SelectOne(&vault, "SELECT * FROM vault WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, vault)
	} else {
		c.JSON(404, gin.H{"error": "vault not found"})
	}

	// curl -i http://localhost:8080/api/v1/vaults/1
}

// PostVault : create new vault
func PostVault(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)

	var vault Vault
	c.Bind(&vault)

	if verbose == true {
		fmt.Println(vault)
		fmt.Println(len(vault.VaultName))
	}

	if len(vault.VaultName) >= 3 { // XXX Check mandatory fields
		err := dbmap.Insert(&vault)
		if err == nil {
			c.JSON(201, vault)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/vaults
}

// UpdateVault : mainly store VerifyKey
func UpdateVault(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	id := c.Params.ByName("id")

	var vault Vault
	err := dbmap.SelectOne(&vault, "SELECT * FROM vault WHERE id=?", id)
	if err == nil {
		var json Vault
		c.Bind(&json)

		if verbose == true {
			fmt.Println(json)
		}

		vaultID, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		vault := Vault{
			Id:        vaultID,
			VerifyKey: json.VerifyKey,
			VaultName: json.VaultName,
			Created:   vault.Created, //vault read from previous select
		}

		if vault.Id != 0 && len(vault.VaultName) >= 3 { // XXX Check mandatory fields
			_, err = dbmap.Update(&vault)
			if err == nil {
				c.JSON(200, vault)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "vault not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/vaults/1
}

// DeleteVault : TODO verify there's no more field with this salt before delete
func DeleteVault(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var vault Vault
	err := dbmap.SelectOne(&vault, "SELECT * FROM vault WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&vault)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "vault not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/vaults/1
}
