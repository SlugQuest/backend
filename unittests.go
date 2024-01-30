package main

// When a new backend function is made, add a test function for it that returns a bool, and then put that func in testmain
import (
	"fmt"
	"strconv"
	"time"
)

func testmain() bool {
	return TestEditTask() && TestGetUserTask() && TestGetTaskId()
}

func TestDeleteTask() bool {
	success, err := DeleteTask(1001)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if !success {
		fmt.Println("something's up")
		return false
	}

	_, _, found := GetTaskId(1001)

	if !found {
		fmt.Println("Delete failed")
		return true
	}

	task := Task{TaskID: 1001, UserID: "1111", Category: "asdf", TaskName: "some name" + strconv.Itoa(1001), Description: "sumdesc" + strconv.Itoa(1001), StartTime: time.Now(), EndTime: time.Now(), IsCompleted: false, IsRecurring: false, IsAllDay: false}
	lol, err := CreateTask(task)
	if lol || (err != nil) {
		return false
	}

	return true
}
func TestEditTask() bool {
	newTask := Task{
		UserID:      "1111",
		TaskID:      3,
		Category:    "example",
		TaskName:    "New Task",
		Description: "Description of the new task",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		IsCompleted: false,
		IsRecurring: false,
		IsAllDay:    false,
	}

	success, err := CreateTask(newTask)
	if err != nil || !success {
		fmt.Println("Error creating task:", err)
		return false
	}

	editedTask := Task{
		TaskID:        3,
		UserID:        "1111",
		Category:      "asdf",
		TaskName:      "edited name",
		Description:   "edited description",
		StartTime:     time.Now(),
		EndTime:       time.Now(),
		IsCompleted:   true,
		IsRecurring:   false,
		IsAllDay:      true,
		RecurringType: "",
		DayOfWeek:     -1,
		DayOfMonth:    -1,
	}

	// Perform the edit
	editSuccess, editErr := EditTask(editedTask, editedTask.TaskID)
	if editErr != nil || !editSuccess {
		fmt.Println("Error editing task:", editErr)
		return false
	}

	taskl, _, _ := GetTaskId(3)
	if taskl.TaskName != "edited name" || !taskl.IsCompleted {
		fmt.Println("Edit verification failed")
		return false
	}

	return true
}
func TestGetUserTask() bool {
	taskl, err := GetUserTask(1111)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if len(taskl) != 500 {
		print("error test get user task wrong count")
		return false
	}
	return true
}
func TestGetTaskId() bool {
	task, erro, found := GetTaskId(1101)
	if erro != nil{
		fmt.Println(erro)
		return false
	}

	if !found {
		fmt.Println("didn't find task")
		return false
	}
	if task.TaskID != 1101{
		fmt.Println("bad task find")
		return false
	}
	task, erro, found = GetTaskId(-5)
	if found{
		fmt.Println("found bad task")
		return false
	}
	return true
}
