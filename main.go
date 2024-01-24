package main

import "github.com/gin-gonic/gin"
import "net/http"

import (
	"fmt"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
  fmt.Println("Running at http://localhost:8080")
	r.Run() // listen and serve on 0.0.0.0:8080
}
