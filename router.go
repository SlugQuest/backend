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
		v1.GET("userTasks/:start/:end", getuserTaskSpan)
		v1.GET("userPoints", getUserPoints)
		v1.GET("getCat/:id", getCategory)
		v1.PUT("makeCat", putCat)
		v1.GET("getBossHealth", getCurrBossHealth)
		v1.GET("/getBoss/:id", getBossById)
		v1.GET("searchuser/:method/:query", searchUsers)
		v1.POST("addFriend/:code", addFriend)
		v1.DELETE("removeFriend/:code", removeFriend)
		v1.GET("user/friends", getFriendList)
		v1.GET("getTeamTask/:id", getTeamTask)
		v1.PUT("addUserTeam/:id/:code", addUserTeam)
		v1.GET("getUserTeams", getUserTeams)
		v1.GET("getTeamUsers/:id", getTeamUsers)
		v1.DELETE("deleteTeamUser/:tid/:code", deleteTeamUser)
		v1.DELETE("deleteTeam/:tid", deleteTeam)
		v1.PUT("createTeam/:name", createTeam)

	}

	return router
}

// Get userID stored in the session
func getUserId(c *gin.Context) (string, error) {
	session := sessions.Default(c)
	userProfile, ok := session.Get("user_profile").(crud.User)
	if !ok {
		log.Printf("getUserId(): error retrieving user id")
		return "", fmt.Errorf("error retrieving user id")
	}
	uid := userProfile.UserID

	return uid, nil
}

func getTeamTask(c *gin.Context) {
	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("getTeamTaskId(): str2int error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
		return
	}

	if tid == crud.NoTeamID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid team"})
		return
	}

	task, err := crud.GetTeamTask(tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve team tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": task})
}

func addUserTeam(c *gin.Context) {
	tid, err := strconv.Atoi(c.Param("id"))
	code := c.Param("code")
	if err != nil {
		log.Println("getTeamTaskId(): str2int error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
		return
	}

	if tid == crud.NoTeamID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid team"})
		return
	}

	succ, err := crud.AddUserToTeam(int64(tid), code)
	if succ && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add user to team"})
	}
}

func getUserTeams(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	ret, err := crud.GetUserTeams(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed teams get"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": ret})
}

func getTeamUsers(c *gin.Context) {
	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("getTeamUsers(): str2int error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
		return
	}

	if tid == crud.NoTeamID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid team"})
		return
	}

	ret, err := crud.GetTeamUsers(int64(tid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed teams get"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": ret})
}

func deleteTeamUser(c *gin.Context) {
	tid, err := strconv.Atoi(c.Param("tid"))
	code := c.Param("code")
	if err != nil {
		log.Println("deleteTeamUser(): str2int error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
		return
	}

	if tid == crud.NoTeamID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid team"})
		return
	}

	ret, err := crud.RemoveUserFromTeam(int64(tid), code)
	if !ret || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error removing user from team"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func deleteTeam(c *gin.Context) {
	tid, err := strconv.Atoi(c.Param("tid"))
	if err != nil {
		log.Println("deleteTeam(): str2int error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This is really bad"})
		return
	}

	if tid == crud.NoTeamID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid team"})
		return
	}

	ret, err := crud.DeleteTeam(int64(tid))
	if !ret || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting team"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func createTeam(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	name := c.Param("name")
	ret, val, err := crud.CreateTeam(name, uid)
	if !ret || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create team"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teamid": val})
}

func passRecurringTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
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

	success, bossId, err := crud.PassRecurringTask(tid, recurrenceID, uid)
	if success && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "bossId": bossId})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pass recurring task"})
		return
	}
}

func failRecurringTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	currBossHealth, err := crud.GetCurrBossHealth(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting current boss health"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"curr_boss_health": currBossHealth})
}

func getCategory(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
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
		log.Printf("getCategory(): unauthorized access to category not owned by current user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Category is not owned by user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Category": cat})
}

func putCat(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
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
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category", "details": err.Error()})
	}
}

func getUserPoints(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	ret, fnd, err := crud.GetUserPoints(uid)
	if !fnd || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving user points"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"points": ret})
}

// Create a new task
func createTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	// Try to bind JSON sent in request into our Task struct
	var json crud.Task
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Printf("createTask(): could not bind request to Task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	json.UserID = uid

	success, taskID, err := crud.CreateTask(json)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "taskID": taskID})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
	}
}

// Edit a task by its ID
func editTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	var json crud.Task
	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Printf("editTask(): could not bind request to Task")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	json.UserID = uid

	// Hall of fame comment
	// fmt.Println("we are checkin the field of the jsawn")
	// fmt.Println(json.Status)

	success, err := crud.EditTask(json, tid)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit task", "details": err.Error()})
	}
}

// Deletes a task by its ID
func deleteTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	arr, err := crud.GetUserTask(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all user tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": arr})
}

func passTheTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	success, bossId, err := crud.PassTask(tid, uid)
	if success && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "bossId": bossId})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pass task"})
	}
}

func failTheTask(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("editTask(): Invalid taskID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TaskId"})
		return
	}

	success, err := crud.FailTask(tid, uid)
	if success && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fail task"})
	}
}

// Retrieve task by ID
func getTaskById(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retreiving task"})
		return
	}
	if task.UserID != uid {
		log.Printf("getTaskById(): unauthorized access to task not owned by current user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "task not owned by current user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// Returns a list of all tasks of the current user
func getuserTaskSpan(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	starttime, err := time.Parse(time.RFC3339, c.Param("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect request time format (start)"})
		return
	}

	endtime, err := time.Parse(time.RFC3339, c.Param("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect request time format (end)"})
		return
	}

	arr, err := crud.GetUserTaskDateTime(uid, starttime, endtime)
	if err != nil {
		log.Println("getAllUserTasks(): Problem probably DB related")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to retreive tasks in given span"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed searching for users"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "format error: accepted methods of user search are by \"name\" or \"code\""})
	}
}

func addFriend(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	their_code := c.Param("code")
	addSuccess, err := crud.AddFriend(uid, their_code)
	if !addSuccess || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func removeFriend(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	their_code := c.Param("code")
	delSuccess, err := crud.DeleteFriend(uid, their_code)
	if !delSuccess || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func getFriendList(c *gin.Context) {
	uid, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating user"})
		return
	}

	friends, err := crud.GetFriendList(uid, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve friend list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"num_friends": len(friends), "list": friends})
}
