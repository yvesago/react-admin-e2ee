package models

/*
  Shared functions for models tests
*/

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func deleteFile(file string) {
	// delete file
	var err = os.Remove(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

type Config struct {
	DBname  string
	Verbose bool
}

func SetConfig(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Verbose", config.Verbose)
		c.Next()
	}
}

// Set test config
var config = Config{
	DBname:  "_test.sqlite3",
	Verbose: true,
}
