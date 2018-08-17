package models

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPeople(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var urla = "/api/v1/peoples"
	router.POST(urla, PostPeople)
	router.GET(urla, GetPeoples)
	router.GET(urla+"/:id", GetPeople)
	router.DELETE(urla+"/:id", DeletePeople)
	router.PUT(urla+"/:id", UpdatePeople)

	b := new(bytes.Buffer)

	router.PUT("/api/v1/utils/:id", UpdateVerifKey)
	log.Println("= http PUT one Util")
	var k = Util{VerifyKey: "XXXXXXXXXXXXXXXX"}
	json.NewEncoder(b).Encode(k)
	req, err := http.NewRequest("PUT", "/api/v1/utils/1", b)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT Salt key success")

	// Add
	log.Println("= http POST People")
	var a = People{Name: "Name test", XAddress: "XXXXXXXXXXXXXXXX address"}
	json.NewEncoder(b).Encode(a)
	req, err = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Add second people
	log.Println("= http POST more People")
	var a2 = People{Name: "Name test2"}
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Test missing mandatory field
	log.Println("= Test missing mandatory field")
	var a2x = People{}
	json.NewEncoder(b).Encode(a2x)
	req, err = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "http POST failed, missing mandatory field")

	// Get all
	log.Println("= http GET all Peoples")
	req, err = http.NewRequest("GET", urla, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []People
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	log.Println("= Test parsing query")
	s := "http://127.0.0.1:8080/api?_filters={\"name\":\"t\"}&_sortDir=ASC&_sortField=created"
	u, _ := url.Parse(s)
	q, _ := url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query, s, _ := ParseQuery(q)
	//fmt.Println(query)
	assert.Equal(t, "name LIKE \"%t%\"", query, "Parse query")
	assert.Equal(t, " ORDER BY datetime(created) ASC", s, "Parse query")

	// Get one
	log.Println("= http GET one People")
	var a1 People
	req, err = http.NewRequest("GET", urla+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.Name, a.Name, "a1 = a")

	// Delete one
	log.Println("= http DELETE one People")
	req, err = http.NewRequest("DELETE", urla+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http DELETE success")
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	req, err = http.NewRequest("GET", urla, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all for count success")
	//fmt.Println(resp.Body)
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 1, len(as), "1 result")

	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "No more /1")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "No more /1")

	// Update one
	log.Println("= http PUT one People")
	a2.Name = "Name test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	var a3 People
	req, err = http.NewRequest("GET", urla+"/2", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one updated success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a2.Name, a3.Name, "a2 Name updated")

	req, _ = http.NewRequest("PUT", urla+"/1", b)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "Can't update missing /1")

	json.NewEncoder(b).Encode(a2x)
	req, _ = http.NewRequest("PUT", urla+"/2", b)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code, "Can't update missing mandatory field in /2")

}

func newRandBytes(length int) ([]byte, error) {
	sb := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, sb); err != nil {
		return nil, err
	}
	return sb, nil
}

func TestE2ECrypto(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var urla = "/api/v1"
	router.POST(urla+"/peoples", PostPeople)
	router.GET(urla+"/peoples/:id", GetPeople)
	router.GET(urla+"/verifkey/:id", GetVerifKey)
	router.PUT(urla+"/verifkey/:id", UpdateVerifKey)
	//router.PUT(urla+"/:id", UpdatePeople)

	log.Println("= Create scrypt key")
	Salt, _ := newRandBytes(12)
	Saltb64 := base64.StdEncoding.EncodeToString(Salt)
	key, _ := scrypt.Key([]byte("secret"), Salt, 16384, 8, 1, 32)

	log.Println("= Create and put Verifykey")
	var aead cipher.AEAD
	aead, _ = chacha20poly1305.New(key)
	Nonce, _ := newRandBytes(12)
	Nonceb64 := base64.StdEncoding.EncodeToString(Nonce)
	Ciphertext := aead.Seal(nil, Nonce[:], []byte("some useless text"), nil)
	Cipherb64 := base64.StdEncoding.EncodeToString(Ciphertext)
	//fmt.Printf("%+v %+v %+v\n", Saltb64, Nonceb64, Cipherb64)

	b := new(bytes.Buffer)
	var v = Util{Id: 1, VerifyKey: Saltb64 + Nonceb64 + Cipherb64}
	json.NewEncoder(b).Encode(v)
	//fmt.Printf("%+v\n", b)
	req, _ := http.NewRequest("PUT", urla+"/verifkey/1", b)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "Verify key set")

	log.Println("= Create People with encrypted XAddress")
	plaintext := "some private people text"
	Nonce, _ = newRandBytes(12)
	Nonceb64 = base64.StdEncoding.EncodeToString(Nonce) // used for authenticate msg
	Ciphertext = aead.Seal(nil, Nonce[:], []byte(plaintext), []byte(Nonceb64))
	Cipherb64 = base64.StdEncoding.EncodeToString(Ciphertext)

	var p = People{Name: "Name test", XAddress: Saltb64 + Nonceb64 + Cipherb64}
	json.NewEncoder(b).Encode(p)
	req, _ = http.NewRequest("POST", urla+"/peoples", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "Add new private data")

	log.Println("= Read People")
	var a1 People
	req, _ = http.NewRequest("GET", urla+"/peoples/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)

	log.Println("= Decrypt XAddress")
	//fmt.Println(a1.XAddress)
	readsalt := a1.XAddress[0:16]
	SaltBin, _ := base64.StdEncoding.DecodeString(readsalt)
	assert.Equal(t, readsalt, Saltb64, "find salt")
	newkey, _ := scrypt.Key([]byte("secret"), SaltBin, 16384, 8, 1, 32)
	aead, _ = chacha20poly1305.New(newkey)
	//fmt.Println(readsalt)
	readnonce := a1.XAddress[16:32]
	readdata := a1.XAddress[32:]
	DataBin, _ := base64.StdEncoding.DecodeString(readdata)
	NonceBin, _ := base64.StdEncoding.DecodeString(readnonce)

	decoded, e := aead.Open(nil, NonceBin, DataBin, []byte(readnonce))
	if e != nil {
		fmt.Printf("1. err %+v\n", e)
	}
	fmt.Println(string(decoded))
	assert.Equal(t, string(decoded), plaintext, "Private text decoded")
}

/*
// Read data from db

func TestCrypto(t *testing.T) {
	dbmap := InitDb("../test.sqlite3")
	var people People
	err := dbmap.SelectOne(&people, "SELECT * FROM people WHERE id=? LIMIT 1", 1)

	if err == nil {
		fmt.Printf("%+v\n", people)
		fmt.Println(people.XAddress)
	}

	salt := people.XAddress[0:16]
	nonce := people.XAddress[16:32]
	data := people.XAddress[32:]

	Salt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}

	Ndecoded, err1 := base64.StdEncoding.DecodeString(nonce)
	if err1 != nil {
		fmt.Println("decode error:", err1)
		return
	}
	//fmt.Println(string(Ndecoded))
	Ddecoded, err2 := base64.StdEncoding.DecodeString(data)
	if err2 != nil {
		fmt.Println("decode error:", err2)
		return
	}
	key, _ := scrypt.Key([]byte("secret"), Salt, 1024, 8, 1, 32)
	fmt.Printf("%X %d\n", key, len(key))
	var aead cipher.AEAD
	if aead, err = chacha20poly1305.New(key); err != nil {
		fmt.Printf("1. err %+v\n", err)
	}
	fmt.Printf("%+v\n", aead.NonceSize())
	fmt.Printf("%+v\n", len(Ddecoded))
	fmt.Printf("%+v\n", Ddecoded)
	if err != nil {
		fmt.Printf("1. err %+v\n", err)
	}
	src := make([]byte, len(Ddecoded))
	res, e := aead.Open(src, Ndecoded, Ddecoded, []byte(nonce))
	if e != nil {
		fmt.Printf("2. err %+v \n", e)
	}
	fmt.Printf("aead: %+v\n Nonce: %+v\n", aead, Ndecoded)
	fmt.Printf("Ddecoded %+v \n", Ddecoded)
	fmt.Printf("res   %+v \n", string(res))
}
*/
