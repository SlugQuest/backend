package unit_tests

import (
	"testing"
	"time"

	"slugquest.com/backend/crud"
	. "slugquest.com/backend/crud"
)

var testTask = Task{
	UserID:         testUser.UserID,
	Category:       1,
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

var recurringTask = Task{
	UserID:         "test_user_id",
	Category:       1,
	TaskName:       "Recurring Test Task",
	Description:    "Sample description",
	StartTime:      time.Now(),
	EndTime:        time.Now().Add(time.Hour),
	Status:         "todo",
	IsRecurring:    true,
	IsAllDay:       false,
	Difficulty:     "easy",
	CronExpression: "0 0 * * *", //every day at midnight
}

func TestGetUserTask(t *testing.T) {
	taskl, err := GetUserTask(testUser.UserID)
	if err != nil {
		t.Errorf("TestGetUserTask(): %v", err)
	}
	if len(taskl) != 500 {
		t.Errorf("TestGetUserTask(): wrong task count, expected 500, got %v", len(taskl))
	}
}

func TestGetUserTaskTime(t *testing.T) {
	starttime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(time.Hour)
	taskl, err := GetUserTaskDateTime(testUser.UserID, starttime, endTime)

	if err != nil {
		t.Errorf("TestGetUserTask(): %v", err)
	}

	if len(taskl) != 500 {
		t.Errorf("TestGetUserTask(): wrong task count, expected 500 got %v", len(taskl))
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

	success, deleteErr := DeleteTask(int(taskID), testTask.UserID)
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
		UserID:         testUser.UserID,
		Category:       1,
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

func TestPassFailTask(t *testing.T) {
	success, taskID, err := CreateTask(testTask)
	if err != nil || !success {
		t.Errorf("TestPassFailTask(): error creating task: %v", err)
	}

	passsucc, err := Passtask(int(taskID), testTask.UserID)
	if err != nil || !passsucc {
		t.Errorf("TestPassFailTask(): error passing task: %v", err)
	}
	task2, _, _ := GetTaskId(int(taskID))
	if task2.Status != "completed" {
		t.Errorf("TestPassFailTask(): wrong status: expected %v, got %v", testTask.Status, task2.Status)
	}

	//points, _, err := GetUserPoints(testUser.UserID)
	failsucc, err := Failtask(int(taskID), testTask.UserID)
	if !failsucc || err != nil {
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

func TestPopRecurringTasksMonth(t *testing.T) {
	success, _, err := CreateTask(recurringTask)
	if err != nil || !success {
		t.Errorf("TestPassFailTask(): error creating task: %v", err)
	}

	popErr := crud.PopRecurringTasksMonth()
	if popErr != nil {
		t.Fatalf("Error populating recurring tasks: %v", popErr)
	}

	count, err := CountRecurringLogEntries()
	if err != nil {
		t.Fatalf("Error counting recurring log entries: %v", err)
	}
	if count != 32 {
		t.Errorf("TestPopRecurringTasksMonth(): wrong count%v", count)
	}
}
