package testing

// When a new backend function is made, add a test function for it that returns a bool, and then put that func in testmain
import (
	"fmt"
	"log"
	"time"

	"path/filepath"

	"slugquest.com/backend/crud"
	. "slugquest.com/backend/crud"
)

var dummyUserID string = "1111"
var testUserID string = "2222" // testing user functions

func RunAllTests() bool {
	ConnectToDB(true)
	dummy_err := LoadDumbData()
	if dummy_err != nil {
		log.Fatalf("error loaduing dumb data: %v", dummy_err)
	}
	return TestGetCurrBossHealth() && TestGetUserTask() && TestGetCategory() && TestDeleteTask() && TestPassFailTask() && TestEditTask() && TestGetTaskId() && TestAddUser() && TestEditUser() && TestDeleteUser()
}

func TestUPoints() bool {
	// NEEDS TO BE DONE
	return false
}

func TestGetCategory() bool {
	cat, bol, erro := GetCatId(50)
	if !bol {
		log.Println("TestGetCat(): Get Cat ID(): cat id not found")
	}
	if erro != nil {
		log.Printf("TestGetCat(): Get Cat ID() #1: %v", erro)
		return false
	}

	if cat.CatID != 50 {
		log.Println("TestGcat(): found wrong cat")
		return false
	}

	cat, bol, erro = GetCatId(-5)
	if bol {
		log.Printf("TestGetCat(): Get Cat ID():  find catad")
		return false
	}
	if erro != nil {
		log.Printf("TestGetCat(): Get Cat ID() #2: %v", erro)
		return false
	}

	return true
}

func TestGetCurrBossHealth() bool {
	newUser := crud.User{
		UserID:   "test_user",
		Username: "test_username",
		Picture:  "test_picture.jpg",
		Points:   10,
		BossId:   1,
	}

	addUserSuccess, addUserErr := crud.AddUser(newUser)
	if addUserErr != nil || !addUserSuccess {
		log.Printf("TestGetCurrBossHealth(): error adding test user: %v", addUserErr)
		return false
	}

	newBoss := crud.Boss{
		BossID: 1,
		Name:   "Test Boss",
		Health: 30,
		Image:  filepath.Join("images", "clown.jpeg"),
	}

	addBossSuccess, addBossErr := AddBoss(newBoss)
	if addBossErr != nil || !addBossSuccess {
		log.Printf("TestGetCurrBossHealth(): error adding test boss: %v", addBossErr)
		return false
	}

	currBossHealth, err := crud.GetCurrBossHealth(newUser.UserID)
	if err != nil {
		log.Printf("TestGetCurrBossHealth(): error getting current boss health: %v", err)
		return false
	}

	if currBossHealth != 20 {
		fmt.Printf("curr boss health: %v", currBossHealth)
		return false
	}
	return true
}

func TestPassFailTask() bool {
	// tx, err := DB.Beginx()

	// // Insert the user into UserTable
	// _, err = tx.Exec("INSERT INTO UserTable (UserID, Points, BossId) VALUES (?, ?, ?)", dummyUserID, 0, 1)
	// if err != nil {
	// 	log.Printf("TestPassFailTask(): error inserting user into UserTable: %v", err)
	// 	return false
	// }

	// tx.Commit()

	newTask := Task{
		UserID:         dummyUserID,
		Category:       "yo",
		TaskName:       "New Task",
		Description:    "Description of the new task",
		StartTime:      time.Now(),
		EndTime:        time.Now().Add(time.Hour),
		Status:         "todo",
		Difficulty:     "hard",
		CronExpression: "",
		IsRecurring:    false,
		IsAllDay:       false,
	}

	success, taskID, err := CreateTask(newTask)
	if err != nil || !success {
		log.Printf("TestPassFailTask(): error creating task: %v", err)
		return false
	}

	passsucc := Passtask(int(taskID))
	if !passsucc {
		log.Printf("TestPassFailTask(): 1 %v", err)
		return false
	}
	task2, _, _ := GetTaskId(int(taskID))
	if task2.Status != "completed" {
		fmt.Printf("TestPassFailTask(): wrong status: %v %v", newTask.Status, task2.Status)
		return false
	}

	//points, _, err := GetUserPoints(dummyUserID)
	failsucc := Failtask(int(taskID))
	if !failsucc {
		log.Printf("TestPassFailTask(): 2 %v", err)
		return false
	}
	// if points != CalculatePoints(newTask.Difficulty) {
	// 	log.Printf("TestPassFailTask(): 3 %v", err)
	// 	return false
	// }

	task3, _, _ := GetTaskId(int(taskID))
	if task3.Status != "failed" {
		log.Printf("TestPassFailTask(): bad value on true fal%v", task3.Status)
		return false
	}
	return true

}

func TestDeleteTask() bool {
	newTask := Task{
		UserID:         dummyUserID,
		Category:       "yo",
		TaskName:       "New Task",
		Description:    "Description of the new task",
		StartTime:      time.Now(),
		EndTime:        time.Now().Add(time.Hour),
		Status:         "completed",
		IsRecurring:    false,
		IsAllDay:       false,
		Difficulty:     "easy",
		CronExpression: "",
	}

	success, taskID, err := CreateTask(newTask)
	if err != nil || !success {
		log.Printf("TestDeleteTask(): error creating task: %v", err)
		return false
	}

	success, deleteErr := DeleteTask(int(taskID))
	if deleteErr != nil {
		log.Printf("TestDeleteTask(): %v", err)
		return false
	}

	if !success {
		log.Println("TestDeleteTask(): something's up")
		return false
	}

	_, bol, _ := GetTaskId(int(taskID))

	if bol {
		log.Println("TestDeleteTask(): delete failed")
		return false
	}

	return true
}
func TestEditTask() bool {
	newTask := Task{
		UserID:         dummyUserID,
		Category:       "yo",
		TaskName:       "New Task",
		Description:    "Description of the new task",
		StartTime:      time.Now(),
		EndTime:        time.Now().Add(time.Hour),
		Status:         "completed",
		IsRecurring:    false,
		IsAllDay:       false,
		Difficulty:     "easy",
		CronExpression: "",
	}

	success, taskID, err := CreateTask(newTask)
	if err != nil || !success {
		log.Printf("TestEditTask(): error creating task: %v", err)
		return false
	}

	editedTask := Task{
		TaskID:         int(taskID),
		UserID:         dummyUserID,
		Category:       "yo",
		TaskName:       "edited name",
		Description:    "edited description",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		Status:         "failed",
		IsRecurring:    false,
		IsAllDay:       true,
		RecurringType:  "",
		Difficulty:     "medium",
		CronExpression: "",
	}

	// Perform the edit
	editSuccess, editErr := EditTask(editedTask, editedTask.TaskID)
	if editErr != nil || !editSuccess {
		log.Printf("TestEditTask(): error editing task: %v", editErr)
		return false
	}

	taskResult, _, _ := GetTaskId(int(taskID))

	if taskResult.TaskName != "edited name" ||
		taskResult.Description != "edited description" ||
		taskResult.Status != "failed" ||
		taskResult.IsAllDay != true ||
		taskResult.Difficulty != "medium" {
		log.Println("TestEditTask(): edit verification failed")
		return false
	}

	return true
}

func TestGetUserTask() bool {
	taskl, err := GetUserTask(dummyUserID)
	if err != nil {
		log.Printf("TestGetUserTask(): %v", err)
		return false
	}
	if len(taskl) != 500 {
		log.Printf("TestGetUserTask(): wrong task count, expected 500 god %v", len(taskl))
		return false
	}
	return true
}

func TestGetUserTaskTime() bool {

	starttime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(time.Hour)
	taskl, err := GetUserTaskDateTime(dummyUserID, starttime, endTime)
	if err != nil {
		log.Printf("TestGetUserTask(): %v", err)
		return false
	}
	if len(taskl) != 500 {
		log.Printf("TestGetUserTask(): wrong task count, expected 500 god %v", len(taskl))
		return false
	}
	return true
}

func TestGetTaskId() bool {
	task, bol, erro := GetTaskId(50)
	if !bol {
		log.Printf("not found")
	}
	if erro != nil {
		log.Printf("TestGetTaskid(): %v", erro)
		return false
	}

	if task.TaskID != 50 {
		log.Println("TestGetTaskId(): found wrong task")
		return false
	}

	task, bol, erro = GetTaskId(-5)
	if bol {
		log.Printf("TestGetTaskid(): find task bad")
		return false
	}
	return true
}
func TestAddUser() bool {
	newUser := User{
		UserID:   testUserID,
		Username: "sluggo",
		Picture:  "lol.jpg",
		Points:   1,
		BossId:   1,
	}

	addSuccess, addErr := AddUser(newUser)
	if addErr != nil || !addSuccess {
		log.Printf("TestAddUser(): couldn't add user")
		return false
	}

	_, found, _ := GetUser(newUser.UserID)
	if !found {
		log.Println("TestAddUser(): add failed")
		return false
	}

	return true
}

func TestEditUser() bool {
	// Original is one inserted in TestAddUser()
	editedUser := User{
		UserID:   testUserID,
		Username: "not in DB, not tested",
		Picture:  "not in DB, not tested",
		Points:   5,
		BossId:   10,
	}

	editSuccess, editErr := EditUser(editedUser, editedUser.UserID)
	if editErr != nil || !editSuccess {
		log.Printf("TestEditUser(): error editing user: %v", editErr)
		return false
	}

	checkE, _, _ := GetUser(editedUser.UserID)
	if checkE.Points != 5 || checkE.BossId != 10 {
		log.Println("TestEditUser(): edit verfication failed")
		return false
	}

	return true
}

func TestDeleteUser() bool {
	deleteSuccess, deleteErr := DeleteUser(testUserID)
	if deleteErr != nil || !deleteSuccess {
		log.Printf("TestDeleteUser(): couldn't delete user")
		return false
	}

	_, found, _ := GetUser(testUserID)
	if found {
		log.Println("TestDeleteUser(): delete failed")
		return false
	}

	return true
}
