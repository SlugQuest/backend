package testing

// When a new backend function is made, add a test function for it that returns a bool, and then put that func in testmain
import (
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
		Category:    5,
		TaskName:    "New Task",
		Description: "Description of the new task",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Status: "failed",
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
	newTask := Task{
		UserID:      testUserId,
		TaskID:      3,
		Category:    5,
		TaskName:    "New Task",
		Description: "Description of the new task",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Status: "completed",
		IsRecurring: false,
		IsAllDay:    false,
		Points: 0,

	}

	success, taskID, err := CreateTask(newTask)
	if err != nil || !success {
		log.Printf("TestEditTask(): error creating task: %v", err)
		return false
	}

	editedTask := Task{
		TaskID:        int(taskID),
		UserID:        testUserId,
		Category:      5,
		TaskName:      "edited name",
		Description:   "edited description",
		StartTime:     time.Now(),
		EndTime:       time.Now(),
		Status:   "failed",
		IsRecurring:   false,
		IsAllDay:      true,
		RecurringType: "",
		DayOfWeek:     -1,
		DayOfMonth:    -1,
	}

	// Perform the edit
	editSuccess, editErr := EditTask(editedTask, editedTask.TaskID)
	if editErr != nil || !editSuccess {
		log.Printf("TestEditTask(): error editing task: %v", editErr)
		return false
	}

	taskl, _, _ := GetTaskId(int(taskID))
	if taskl.TaskName != "edited name" {
		log.Println("TestEditTask(): edit verfication failed")
		return false
	}

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
