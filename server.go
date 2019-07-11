package main

import (
	. "./models"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"time"
)

// SetConfig gin Middlware to push some config values
func SetConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("CorsOrigin", "*")
		c.Set("Verbose", true)
		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(Database("test.sqlite3"))
	r.Use(SetConfig())
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, Bearer",
		ExposedHeaders:  "x-total-count, Content-Range",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	v1 := r.Group("api/v1")
	{
		v1.GET("/users", GetUsers)
		v1.GET("/users/:id", GetUser)
		v1.POST("/users", PostUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
		v1.OPTIONS("/users", Options)     // POST
		v1.OPTIONS("/users/:id", Options) // PUT, DELETE

		v1.GET("/peoples", GetPeoples)
		v1.GET("/peoples/:id", GetPeople)
		v1.POST("/peoples", PostPeople)
		v1.PUT("/peoples/:id", UpdatePeople)
		v1.DELETE("/peoples/:id", DeletePeople)
		v1.OPTIONS("/peoples", Options)     // POST
		v1.OPTIONS("/peoples/:id", Options) // PUT, DELETE

		v1.GET("/vault", GetVaults)
		v1.GET("/vault/:id", GetVault)
		v1.POST("/vault", PostVault)
		v1.PUT("/vault/:id", UpdateVault)
		v1.DELETE("/vault/:id", DeleteVault)
		v1.OPTIONS("/vault", Options)     // POST
		v1.OPTIONS("/vault/:id", Options) // PUT, DELETE

	}

	r.Run("localhost:8088")
}

// Options common response for rest options
func Options(c *gin.Context) {
	Origin := c.MustGet("CorsOrigin").(string)

	c.Writer.Header().Set("Access-Control-Allow-Origin", Origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}
