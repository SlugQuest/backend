package main

import (
	"net/http"
	"strconv"

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
	var json Task //instance of Task struct defined in handler
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := EditTask(json, id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit task", "details": err.Error()})
	}
}

func deleteTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Called deleteTask"})
}

func getAllUserTasks(c *gin.Context) {

	uid := 1111
	arr, err := GetUserTask(uid)
	if err != nil {
		fmt.Println("ERROR LOG:  Problem in getAllUserTasks, probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
	}
	c.JSON(http.StatusOK, gin.H{"list": arr})
}

func getTaskById(c *gin.Context) {
	tid, err1 := strconv.Atoi(c.Param("id"))
	if(err1 != nil){
		fmt.Println("ERROR LOG:  str2int error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
	}
	task, err := GetTaskId(tid)
	if err != nil {
		fmt.Println("ERROR LOG:  Problem in getAllUserTasks, probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}
