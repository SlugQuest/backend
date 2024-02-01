// As provided by Auth0

package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"slugquest.com/backend/authentication"
	"slugquest.com/backend/crud"
)

// New registers the routes and returns the router.
func CreateRouter(auth *authentication.Authenticator) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	// Set up cookie store for the user session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// Router: takes incoming requests and routes them to functions to handle them
	router.GET("/login", authentication.LoginHandler(auth))
	router.GET("/logout", authentication.LogoutHandler)
	router.GET("/callback", authentication.CallbackHandler(auth))
	// router.GET("/user", authentication.UserProfileHandler)

	// Building a group of routes starting with this path
	v1 := router.Group("/main/blah") //TODO: FIX the route and the uri's below
	{
		// First middleware to use is verifying authentication
		v1.Use(authentication.IsAuthenticated)

		v1.GET("tasks", getAllUserTasks)
		v1.GET("task/:id", getTaskById)
		v1.POST("tasks", createTask)
		v1.PUT("tasks/:id", editTask)
		v1.DELETE("tasks/:id", deleteTask)
	}

	return router
}

// Create a new task
func createTask(c *gin.Context) {
	var json crud.Task //instance of Task struct defined in db_handler.go

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} //take any JSON sent in the BODY of the request and try to bind it to our Task struct

	success, taskID, err := crud.CreateTask(json) //pass struct into function to add Task to db
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "taskID": taskID})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}
}

// Edit a task by its ID
func editTask(c *gin.Context) {
	var json crud.Task //instance of Task struct defined in handler

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := crud.EditTask(json, id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit task", "details": err.Error()})
		return
	}
}

// Deletes a task by its ID
func deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("deleteTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	success, err := crud.DeleteTask(id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task", "details": err.Error()})
	}
}

// Returns a list of all tasks of the current user
func getAllUserTasks(c *gin.Context) {
	// Retrieve current session variables and cookies
	session := sessions.Default(c)

	// user_id stored as a variable within the session
	uid := session.Get("user_id").(string)
	arr, err := crud.GetUserTask(uid)
	if err != nil {
		log.Println("getAllUserTasks(): Problem probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": arr})
}

// Retrieve task by ID
func getTaskById(c *gin.Context) {
	tid, err1 := strconv.Atoi(c.Param("id"))
	if err1 != nil {
		log.Println("getTaskById(): str2int error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}

	task, value, err := crud.GetTaskId(tid)
	if !value {
		log.Printf("getTaskById(): Did not find task with ID %v", tid)
		c.JSON(http.StatusBadRequest, gin.H{"not found": "no task"})
		return
	}
	if err != nil {
		log.Println("getTaskById(): Problem in getAllUserTasks, probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}
