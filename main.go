package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"slugquest.com/backend/middleware"

	"fmt"

	_ "github.com/mattn/go-sqlite3"

	envfuncs "github.com/joho/godotenv"
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

	err := connectToDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	erro := loadDumbData()
	if erro != nil {
		fmt.Println("error loaduing dumb data", err)
	}
	utest := testmain()
	if !utest {
		fmt.Println("unit test failure")
		return
	}

	// Router: takes incoming requests and routes them to functions to handle them
	// Building a group of routes starting with this path

	// This would make it a global handler to pass through each route, but I assume there will
	// be public endpoints later, leaving it commented out for now
	// r.Use(middleware.EnsureValidToken())

	v1 := r.Group("/main/blah") //TODO: FIX the route and the uri's below
	{
		v1.GET("tasks", middleware.EnsureValidToken(), getAllUserTasks)
		v1.GET("task/:id", middleware.EnsureValidToken(), getTaskById)
		v1.POST("tasks", middleware.EnsureValidToken(), createTask)
		v1.PUT("tasks/:id", middleware.EnsureValidToken(), editTask)
		v1.DELETE("tasks/:id", middleware.EnsureValidToken(), deleteTask)

	}

	log.Print("Running at http://localhost:8080")
	r.Run() // listen and serve on 0.0.0.0:8080
}

func createTask(c *gin.Context) {

	var json Task //instance of Task struct defined in handler

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} //take any JSON sent in the BODY of the request and try to bind it to our Task struct

	success, err, taskID := CreateTask(json) //pass struct into function to add Task to db

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "taskID": taskID})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}
}

func editTask(c *gin.Context) {
	var json Task //instance of Task struct defined in handler
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := EditTask(json, id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit task", "details": err.Error()})
		return
	}
}

func deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
	}

	success, err := DeleteTask(id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task", "details": err.Error()})
	}
}

func getAllUserTasks(c *gin.Context) {

	uid := 1111
	arr, err := GetUserTask(uid)
	if err != nil {
		fmt.Println("ERROR LOG:  Problem in getAllUserTasks, probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": arr})
}

func getTaskById(c *gin.Context) {
	tid, err1 := strconv.Atoi(c.Param("id"))
	if err1 != nil {
		fmt.Println("ERROR LOG:  str2int error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}
	task, err, value := GetTaskId(tid)
	if !value {
		fmt.Println("ERROR LOG:  getting a non idd task")
		c.JSON(http.StatusBadRequest, gin.H{"not found": "no task"})
		return
	}
	if err != nil {
		fmt.Println("ERROR LOG:  Problem in getAllUserTasks, probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}
