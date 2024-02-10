// As provided by Auth0

package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"strconv"
	"time"

	
	"github.com/gin-contrib/cors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"slugquest.com/backend/authentication"
	"slugquest.com/backend/crud"
)

// New registers the routes and returns the router.
func CreateRouter(auth *authentication.Authenticator) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://" + authentication.FRONTEND_HOST},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})
	gob.Register(crud.User{})
	// Set up cookie store for the user session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// Router: takes incoming requests and routes them to functions to handle them
	router.GET("/login", authentication.LoginHandler(auth))
	router.GET("/logout", authentication.LogoutHandler)
	router.GET("/callback", authentication.CallbackHandler(auth))
	// router.GET("/user", authentication.UserProfileHandler)

	// Building a group of routes starting with this path
	v1 := router.Group("/api/v1") //TODO: FIX the route and the uri's below
	{
		// First middleware to use is verifying authentication
		v1.Use(authentication.IsAuthenticated)

		v1.GET("tasks", getAllUserTasks)
		v1.GET("task/:id", getTaskById)
		v1.POST("task", createTask)
		v1.POST("passtask/:id", passTheTask)
		v1.POST("failtask/:id", failTheTask)
		v1.PUT("task/:id", editTask)
		v1.DELETE("task/:id", deleteTask)
		v1.GET("userTasks/:id/:start/:end", getuserTaskSpan)
		v1.GET("userPoints", getUserPoints)
		v1.GET("getCat/:id", getCategory)
		v1.PUT("makeCat", putCat)
	}

	return router
}

func getCategory(c *gin.Context) {

	cid, err1 := strconv.Atoi(c.Param("id"))
	if err1 != nil {
		log.Println("getCategory): str2int error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}

	Cat, bol, err := crud.GetCatId(cid)
	if !bol {
		log.Println("getTaskById(): Problem in getAllUserTasks, probably DB related", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Category": Cat})
}

func putCat(c *gin.Context) {
	var json crud.Category //instance of Task struct defined in db_handler.go

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} //take any JSON sent in the BODY of the request and try to bind it to our Task struct

	success, catID, err := crud.CreateCategory(json) //pass struct into function to add Task to db
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "catID": catID})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cat", "details": err.Error()})
		return
	}
}
func getUserPoints(c *gin.Context) {
	//PLACEHOLDER VALUE
	uid := "1111"
	ret, fnd, err := crud.GetUserPoints(uid)

	if !fnd {
		log.Println("getTaskById(): Problem in getUserPoints, probably DB related", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"points": ret})
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
	// Retrieve the user_id through the struct stored in the session
	session := sessions.Default(c)
	userProfile, ok := session.Get("user_profile").(crud.User)
	if !ok {
		c.String(http.StatusInternalServerError, "Couldn't retreive user's id to display tasks.")
		return
	}
	uid := userProfile.UserID
	arr, err := crud.GetUserTask(uid)
	if err != nil {
		log.Println("getAllUserTasks(): Problem probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": arr})
}

func passTheTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}
	erro := crud.Passtask(id)

	if !erro {

		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	} else {
		log.Println(erro)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pass task"})
		return
	}
}

func failTheTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	erro := crud.Failtask(id)

	if !erro {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	} else {
		log.Println(erro)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fail task"})
		return
	}
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
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// Returns a list of all tasks of the current user
func getuserTaskSpan(c *gin.Context) {
	// TODO: ill be fixing this
	// user_id stored as a variable within the session
	// uid := c.GetString("user_id")
	// log.Printf("found userid = %v", uid)
	// if uid == "" {
	// 	log.Println("getAllUserTasks(): couldn't get user_id")
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive user id"})
	// 	return
	// }
	uid := authentication.Curr_user_id
	starttime, err1 := time.Parse(time.RFC3339, c.GetString("start"))
	if err1 != nil {
		log.Println("Please pass in a well formatted time. This is a frontend issue.")
		return
	}

	endtime, err2 := time.Parse(time.RFC3339, c.GetString("end"))
	if err2 != nil {
		log.Println("Please pass in a well formatted time. This is a frontend issue.")
		return
	}

	arr, err := crud.GetUserTaskDateTime(uid, starttime, endtime)
	if err != nil {
		log.Println("getAllUserTasks(): Problem probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is really bad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": arr})
}
