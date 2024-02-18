package testing

import (
	"fmt"
	"log"
	"time"

	. "slugquest.com/backend/crud"
)

var userForTaskTable = User{
	UserID:   "tasktable_user_id",
	Username: "sluggo1",
	Picture:  "lol.jpg",
	Points:   1,
	BossId:   1,
}

var testTask = Task{
	UserID:         userForTaskTable.UserID,
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

func TestPassFailTask() bool {
	// tx, err := DB.Beginx()

	// // Insert the user into UserTable
	// _, err = tx.Exec("INSERT INTO UserTable (UserID, Points, BossId) VALUES (?, ?, ?)", userForTaskTable.UserID, 0, 1)
	// if err != nil {
	// 	log.Printf("TestPassFailTask(): error inserting user into UserTable: %v", err)
	// 	return false
	// }

	// tx.Commit()

	success, taskID, err := CreateTask(testTask)
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
		fmt.Printf("TestPassFailTask(): wrong status: %v %v", testTask.Status, task2.Status)
		return false
	}

	//points, _, err := GetUserPoints(userForTaskTable.UserID)
	failsucc := Failtask(int(taskID))
	if !failsucc {
		log.Printf("TestPassFailTask(): 2 %v", err)
		return false
	}
	// if points != CalculatePoints(testTask.Difficulty) {
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

func TestGetUserTask() bool {
	taskl, err := GetUserTask(userForTaskTable.UserID)
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
	taskl, err := GetUserTaskDateTime(userForTaskTable.UserID, starttime, endTime)
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

func TestDeleteTask() bool {
	success, taskID, err := CreateTask(testTask)
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
	success, taskID, err := CreateTask(testTask)
	if err != nil || !success {
		log.Printf("TestEditTask(): error creating task: %v", err)
		return false
	}

	editedTask := Task{
		TaskID:         int(taskID),
		UserID:         userForTaskTable.UserID,
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
