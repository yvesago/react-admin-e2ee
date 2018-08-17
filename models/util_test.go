package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	//	"net/url"
	"testing"
)

func TestVerifyKey(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var urla = "/api/v1/utils"
	router.GET(urla+"/:id", GetVerifKey)
	router.PUT(urla+"/:id", UpdateVerifKey)

	var a = Util{Id: 1, VerifyKey: ""}
	var a2 = Util{Id: 1, VerifyKey: "VerifyKey test2"}
	b := new(bytes.Buffer)

	// Get one
	log.Println("= http GET VerifyKey")
	var a1 Util
	req, err := http.NewRequest("GET", urla+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.VerifyKey)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.VerifyKey, a.VerifyKey, "a1 = a")

	// Update VerifyKey
	log.Println("= http PUT VerifyKey")
	//var a4 = Util{VerifyKey: "VerifyKey test2 updated"}
	a2.VerifyKey = "VerifyKey test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	var a3 Util
	req, err = http.NewRequest("GET", urla+"/2", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all updated success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	//fmt.Println(a1.VerifyKey)
	//fmt.Println(resp.Body)
	assert.Equal(t, a2.VerifyKey, a3.VerifyKey, "a2 VerifyKey updated")

}
