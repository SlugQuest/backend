package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
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
	if erro != nil {
		fmt.Println("breaky")
	}
	fmt.Println(string(schemacreate))
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Println("breaky")
	}

	// force a connection and test that it worked
	swagmoney := db.Ping()

	db.MustExec(string(schemacreate))
	if swagmoney != nil {
		fmt.Println("breaky")
	} else {
		fmt.Println("not breaky")
	}

	//Router: takes incoming requests and routes them to functions to handle them
	//Building a group of routes starting with this path
	v1 := r.Group("/main/blah") //TODO: FIX the route and the uri's below
	{
		v1.GET("tasks", getAllUserTasks)
		v1.GET("task/:id", getTaskById)
		v1.POST("tasks", createTask)
		v1.PUT("tasks/:id", editTask)
		v1.DELETE("tasks/:id", deleteTask)

	}

	fmt.Println("Running at http://localhost:8080")
	r.Run() // listen and serve on 0.0.0.0:8080
}

func createTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Called createTask"})
}

func editTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Called editTask"})
}

func deleteTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Called deleteTask"})
}

func getAllUserTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Called getAllUserTasks"})
}

func getTaskById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Called GetTaskById on Id:" + id})
}
