package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

	err := connectToDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
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

	var json Task //instance of Task struct defined in handler

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} //take any JSON sent in the BODY of the request and try to bind it to our Task struct

	success, err := CreateTask(json) //pass struct into function to add Task to db

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
	}
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
