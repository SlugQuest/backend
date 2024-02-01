package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	envfuncs "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"slugquest.com/backend/crud"
	"slugquest.com/backend/testing"
)

func main() {
	// Load .env
	if err := envfuncs.Load(); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := crud.ConnectToDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	erro := crud.LoadDumbData()
	if erro != nil {
		fmt.Println("error loaduing dumb data", err)
	}
	utest := testing.RunAllTests()
	if !utest {
		fmt.Println("unit test failure")
		return
	}

	log.Print("Running at http://localhost:8080")
	r.Run() // listen and serve on 0.0.0.0:8080
}
