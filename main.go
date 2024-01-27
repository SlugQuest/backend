package main

import "github.com/gin-gonic/gin"
import "net/http"
import "os"
import "github.com/jmoiron/sqlx"
import _"github.com/mattn/go-sqlite3"

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
    // exactly the same as the built-in
	schemacreate, erro := os.ReadFile("schema.sql")
	if(erro!=nil){
		fmt.Println("breaky")
	}
	fmt.Println(string(schemacreate))
    db,err := sqlx.Open("sqlite3", ":memory:")
	if err != nil{
		fmt.Println("breaky")
	}
     
    // force a connection and test that it worked
    swagmoney := db.Ping()

	db.MustExec(string(schemacreate))
	if swagmoney != nil{
		fmt.Println("breaky")
	}else{
		fmt.Println("not breaky")
	}
     
  fmt.Println("Running at http://localhost:8080")
	r.Run() // listen and serve on 0.0.0.0:8080
}
