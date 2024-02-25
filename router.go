package main

import (
	"encoding/gob"
	"fmt"
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
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length"},
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
	router.GET("/login", authentication.LoginHandler(auth, false))
	router.GET("/logout", authentication.LogoutHandler)
	router.GET("/callback", authentication.CallbackHandler(auth))
	router.GET("/signup", authentication.LoginHandler(auth, true))

	// Building a group of routes starting with this path
	v1 := router.Group("/api/v1")
	{
		// Verifying authenticated before any of the endpoints for this group
		v1.Use(authentication.IsAuthenticated)

		v1.GET("user", authentication.UserProfileHandler)

		v1.GET("tasks", getAllUserTasks)
		v1.GET("task/:id", getTaskById)
		v1.POST("task", createTask)
		v1.POST("passtask/:id", passTheTask)
		v1.POST("failtask/:id", failTheTask)
		v1.POST("passRecurringTask/:id/:recurrenceID", passRecurringTask)
		v1.POST("failRecurringTask/:id/:recurrenceID", failRecurringTask)
		v1.PUT("task/:id", editTask)
		v1.DELETE("task/:id", deleteTask)
		v1.GET("userTasks/:id/:start/:end", getuserTaskSpan)
		v1.GET("userPoints", getUserPoints)
		v1.GET("getCat/:id", getCategory)
		v1.PUT("makeCat", putCat)
		v1.GET("getBossHealth", getCurrBossHealth)
		v1.GET("/getBoss/:id", getBossById)
		v1.GET("searchuser/:method/:query", searchUsers)
		v1.POST("addFriend/:code", addFriend)
		v1.DELETE("removeFriend/:code", removeFriend)
		v1.GET("user/friends", getFriendList)
	}

	return router
}

// Get userID stored in the session
func getUserId(c *gin.Context) (string, error) {
	session := sessions.Default(c)
	userProfile, ok := session.Get("user_profile").(crud.User)
	if !ok {
		return "", fmt.Errorf("couldn't get user id")
	}
	uid := userProfile.UserID

	return uid, nil
}

func passRecurringTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("passRecurringTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	recurrenceID, err := strconv.Atoi(c.Param("recurrenceID"))
	if err != nil {
		log.Println("passRecurringTask(): Invalid recurrenceID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RecurrenceID"})
		return
	}

	success, err := crud.PassRecurringTask(tid, recurrenceID, uid)
	if success && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	} else {
		log.Printf("passRecurringTask(): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pass recurring task"})
		return
	}
}

func failRecurringTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("failRecurringTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	recurrenceID, err := strconv.Atoi(c.Param("recurrenceID"))
	if err != nil {
		log.Println("failRecurringTask(): Invalid recurrenceID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RecurrenceID"})
		return
	}

	success, err := crud.FailRecurringTask(tid, recurrenceID, uid)
	if success && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	} else {
		log.Printf("failRecurringTask(): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fail recurring task"})
		return
	}
}

func getBossById(c *gin.Context) {
	bossID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BossID"})
		return
	}

	boss, exists, err := crud.GetBossById(bossID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve boss"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Boss not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"boss": boss})
}

func getCurrBossHealth(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	currBossHealth, err := crud.GetCurrBossHealth(uid)
	if err != nil {
		c.String(http.StatusInternalServerError, "Can't find boss health")
		return
	}
	c.JSON(http.StatusOK, gin.H{"curr_boss_health": currBossHealth})
}

func getCategory(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("getCategory(): str2int error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "uh oh"})
		return
	}

	cat, found, err := crud.GetCatId(cid)
	if err != nil {
		log.Printf("getCategory(): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find category"})
		return
	}
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not find category"})
		return
	}

	if cat.UserID != uid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Category is not owned by user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Category": cat})
}

func putCat(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	var json crud.Category
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	json.UserID = uid

	success, catID, err := crud.CreateCategory(json)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "catID": catID})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cat", "details": err.Error()})
		return
	}
}

func getUserPoints(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

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
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	var json crud.Task //instance of Task struct defined in db_handler.go
	if err := c.ShouldBindJSON(&json); err != nil {
		fmt.Println("errorcasexsit", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} //take any JSON sent in the BODY of the request and try to bind it to our Task struct
	json.UserID = uid
	fmt.Println("creat")
	fmt.Println(json)

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
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	var json crud.Task //instance of Task struct defined in handler

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	json.UserID = uid

	success, err := crud.EditTask(json, tid)

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
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("deleteTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	success, err := crud.DeleteTask(tid, uid)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task", "details": err.Error()})
	}
}

// Returns a list of all tasks of the current user
func getAllUserTasks(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	arr, err := crud.GetUserTask(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all user tasks"})
		return
	}
	log.Println("working")
	log.Println(arr)
	c.JSON(http.StatusOK, gin.H{"list": arr})
}

func passTheTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	success, err := crud.Passtask(tid, uid)
	if success && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	} else {
		log.Printf("passTheTask(): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pass task"})
		return
	}
}

func failTheTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	success, err := crud.Failtask(tid, uid)
	if success && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	} else {
		log.Printf("failTheTask(): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fail task"})
		return
	}
}

// Retrieve task by ID
func getTaskById(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	tid, err1 := strconv.Atoi(c.Param("id"))
	if err1 != nil {
		log.Println("getTaskById(): str2int error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
		return
	}

	task, value, err := crud.GetTaskId(tid)
	if !value {
		log.Printf("getTaskById(): Did not find task with ID %v", tid)
		c.JSON(http.StatusBadRequest, gin.H{"not found": "no task"})
		return
	}
	if err != nil {
		log.Printf("getTaskById(): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
		return
	}
	if task.UserID != uid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "task not owned by user"})
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// Returns a list of all tasks of the current user
func getuserTaskSpan(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	starttime, err1 := time.Parse(time.RFC3339, c.GetString("start"))
	if err1 != nil {
		c.String(http.StatusBadRequest, "Error: incorrect request time format")
		return
	}

	endtime, err2 := time.Parse(time.RFC3339, c.GetString("end"))
	if err2 != nil {
		c.String(http.StatusBadRequest, "Error: incorrect request time format")
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

func searchUsers(c *gin.Context) {
	query := c.Param("query")
	method := c.Param("method")
	if method == "name" {
		users, foundAny, err := crud.SearchUsername(query, false)
		if err != nil {
			log.Printf("searchUsers(): error searching for users: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
			return
		}

		if !foundAny {
			c.JSON(http.StatusNotFound, gin.H{"message": "No users found", "num_results": 0, "users": users})
			return
		}

		c.JSON(http.StatusOK, gin.H{"num_results": len(users), "users": users})
	} else if method == "code" {
		res := []map[string]interface{}{}
		user, foundOne, err := crud.SearchUserCode(query, false)
		if err != nil {
			log.Printf("searchUsers(): error searching for users: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
			return
		}

		if !foundOne {
			c.JSON(http.StatusNotFound, gin.H{"message": "No user with this social code found", "num_results": 0, "users": res})
			return
		}

		res = append(res, user)
		c.JSON(http.StatusOK, gin.H{"num_results": 1, "users": res})

	} else {
		c.String(http.StatusBadRequest, "Format error: accepted methods of user search are by \"name\" and \"code\".")
	}
}

func addFriend(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	their_code := c.Param("code")
	addSuccess, err := crud.AddFriend(uid, their_code)
	if !addSuccess || err != nil {
		log.Printf("addFriend(): error adding friend: %v", err)
		c.String(http.StatusInternalServerError, "Failure to add friend")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func removeFriend(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	their_code := c.Param("code")
	delSuccess, err := crud.DeleteFriend(uid, their_code)
	if !delSuccess || err != nil {
		log.Printf("deleteFriend(): error removing friend: %v", err)
		c.String(http.StatusInternalServerError, "Failure to remove friend")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func getFriendList(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failure to retrieve user id")
		return
	}

	friends, err := crud.GetFriendList(uid, false)
	if err != nil {
		log.Printf("getFriendList(): error retrieving friend list: %v", err)
		c.String(http.StatusInternalServerError, "Failure to retrieve friend list")
		return
	}

	c.JSON(http.StatusOK, gin.H{"num_friends": len(friends), "list": friends})
}
