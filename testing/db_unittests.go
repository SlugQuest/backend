package testing

// When a new backend function is made, add a test function for it that returns a bool, and then put that func in testmain
import (
	"fmt"
	"log"
	"time"

	. "slugquest.com/backend/crud"
)

var testUserId string = "1111"

func RunAllTests() bool {
	ConnectToDB(true)
	dummy_err := LoadDumbData()
	if dummy_err != nil {
		log.Fatalf("error loaduing dumb data: %v", dummy_err)
	}
	return TestGetUserTask() && TestDeleteTask() && TestEditTask() && TestGetTaskId()
}

func TestDeleteTask() bool {
	newTask := Task{
		UserID:      testUserId,
		Category:    "yo",
		TaskName:    "New Task",
		Description: "Description of the new task",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Status:      "failed",
		IsRecurring: false,
		IsAllDay:    false,
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

	_, found, _ := GetTaskId(int(taskID))

	if found {
		log.Println("TestDeleteTask(): delete failed")
		return false
	}

	return true
}
func TestEditTask() bool {
	tx, err := DB.Beginx()
	newTask := Task{
		UserID:         testUserId,
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

	var userExists bool
	err = tx.Get(&userExists, "SELECT EXISTS (SELECT 1 FROM UserTable WHERE UserID = ?)", newTask.UserID)
	if err != nil {
		fmt.Println("CreateTask(): breaky 2", err)
		return false
	}

	if !userExists {
		_, err = tx.Exec("INSERT INTO UserTable (UserID, Points) VALUES (?, 0)", newTask.UserID)
		if err != nil {
			fmt.Println("CreateTask(): breaky 3", err)
			return false
		}
	}

	success, taskID, err := CreateTask(newTask)
	if err != nil || !success {
		log.Printf("TestEditTask(): error creating task: %v", err)
		return false
	}

	editedTask := Task{
		TaskID:         int(taskID),
		UserID:         testUserId,
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

	taskResult, found, _ := GetTaskId(int(taskID))
	if !found {
		log.Println("TestEditTask(): edited task not found")
		return false
	}

	if taskResult.TaskName != "edited name" ||
		taskResult.Description != "edited description" ||
		taskResult.Status != "failed" ||
		taskResult.IsAllDay != true ||
		taskResult.Difficulty != "medium" ||
		taskResult.CronExpression != "newcron" {
		log.Println("TestEditTask(): edit verification failed")
		return false
	}

	//newPoints := CalculatePoints("medium")

	// // user, _, _ := GetUserById(testUserId)
	// if user.Points != newPoints {
	// 	log.Println("TestEditTask(): user points verification failed")
	// 	return false
	// }

	return true
}

func TestGetUserTask() bool {
	taskl, err := GetUserTask(testUserId)
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
	task, found, erro := GetTaskId(50)
	if erro != nil {
		log.Printf("TestGetTaskid(): %v", erro)
		return false
	}

	if !found {
		log.Println("TestGetTaskId(): didn't find task")
		return false
	}
	if task.TaskID != 50 {
		log.Println("TestGetTaskId(): found wrong task")
		return false
	}

	task, found, erro = GetTaskId(-5)
	if erro != nil {
		log.Printf("TestGetTaskid(): %v", erro)
		return false
	}
	if found {
		log.Println("TestGetTaskId(): found task by invalid id")
		return false
	}
	return true
}
