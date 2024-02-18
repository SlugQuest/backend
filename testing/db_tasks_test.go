package testing

import (
	"testing"
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

func TestPassFailTask(t *testing.T) {
	// tx, err := DB.Beginx()

	// // Insert the user into UserTable
	// _, err = tx.Exec("INSERT INTO UserTable (UserID, Points, BossId) VALUES (?, ?, ?)", userForTaskTable.UserID, 0, 1)
	// if err != nil {
	// 	t.Errorf("TestPassFailTask(): error inserting user into UserTable: %v", err)
	// 	return false
	// }

	// tx.Commit()

	success, taskID, err := CreateTask(testTask)
	if err != nil || !success {
		t.Errorf("TestPassFailTask(): error creating task: %v", err)
	}

	passsucc := Passtask(int(taskID))
	if !passsucc {
		t.Errorf("TestPassFailTask(): 1 %v", err)
	}
	task2, _, _ := GetTaskId(int(taskID))
	if task2.Status != "completed" {
		t.Errorf("TestPassFailTask(): wrong status: %v %v", testTask.Status, task2.Status)
	}

	//points, _, err := GetUserPoints(userForTaskTable.UserID)
	failsucc := Failtask(int(taskID))
	if !failsucc {
		t.Errorf("TestPassFailTask(): 2 %v", err)
	}
	// if points != CalculatePoints(testTask.Difficulty) {
	// 	t.Errorf("TestPassFailTask(): 3 %v", err)
	// 	return false
	// }

	task3, _, _ := GetTaskId(int(taskID))
	if task3.Status != "failed" {
		t.Errorf("TestPassFailTask(): bad value on true fal%v", task3.Status)
	}
}

func TestGetUserTask(t *testing.T) {
	taskl, err := GetUserTask(userForTaskTable.UserID)
	if err != nil {
		t.Errorf("TestGetUserTask(): %v", err)
	}
	if len(taskl) != 500 {
		t.Errorf("TestGetUserTask(): wrong task count, expected 500 god %v", len(taskl))
	}
}

func TestGetUserTaskTime(t *testing.T) {
	starttime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(time.Hour)
	taskl, err := GetUserTaskDateTime(userForTaskTable.UserID, starttime, endTime)

	if err != nil {
		t.Errorf("TestGetUserTask(): %v", err)
	}

	if len(taskl) != 500 {
		t.Errorf("TestGetUserTask(): wrong task count, expected 500 god %v", len(taskl))
	}
}

func TestGetTaskId(t *testing.T) {
	task, bol, erro := GetTaskId(50)
	if !bol {
		t.Errorf("TestGetTaskid(): not found")
	}
	if erro != nil {
		t.Errorf("TestGetTaskid(): %v", erro)
	}

	if task.TaskID != 50 {
		t.Error("TestGetTaskId(): found wrong task")
	}

	task, bol, erro = GetTaskId(-5)
	if bol || erro != nil {
		t.Errorf("TestGetTaskid(): find task bad")
	}
}

func TestDeleteTask(t *testing.T) {
	success, taskID, err := CreateTask(testTask)
	if err != nil || !success {
		t.Errorf("TestDeleteTask(): error creating task: %v", err)
	}

	success, deleteErr := DeleteTask(int(taskID))
	if deleteErr != nil {
		t.Errorf("TestDeleteTask(): %v", err)
	}

	if !success {
		t.Error("TestDeleteTask(): something's up")
	}

	_, bol, _ := GetTaskId(int(taskID))

	if bol {
		t.Error("TestDeleteTask(): delete failed")
	}
}

func TestEditTask(t *testing.T) {
	success, taskID, err := CreateTask(testTask)
	if err != nil || !success {
		t.Errorf("TestEditTask(): error creating task: %v", err)
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
		t.Errorf("TestEditTask(): error editing task: %v", editErr)
	}

	taskResult, _, _ := GetTaskId(int(taskID))

	if taskResult.TaskName != "edited name" ||
		taskResult.Description != "edited description" ||
		taskResult.Status != "failed" ||
		taskResult.IsAllDay != true ||
		taskResult.Difficulty != "medium" {
		t.Error("TestEditTask(): edit verification failed")
	}

}
