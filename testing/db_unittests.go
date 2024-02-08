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
	return TestGetUserTask() && TestDeleteTask() && TestPassFailTask() &&/TestEditTask()&&  TestGetTaskId()
}

func TestPassFailTask() bool{
	newTask := Task{
		UserID:      testUserId,
		Category:    "yo",
		TaskName:    "New Task",
		Description: "Description of the new task",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Status:      "completed",
		Difficulty: "hard",
		CronExpression: "asdf",
		IsRecurring: false,
		IsAllDay:    false,
	}

	success, taskID, err := CreateTask(newTask)
	if err != nil || !success {
		log.Printf("TestPassFailTask(): error creating task: %v", err)
		return false
	}

	passsucc := Passtask(int(taskID))
	if !passsucc {
		log.Printf("TestPassFailTask(): 1 %v")
		return false
	}
	task2, _, _:= GetTaskId(int(taskID))
	if task2.Status != "completed" {
		log.Printf("TestPassFailTask(): bad value on true fal%v", task2.Status)
		return false
	}
	failsucc := Failtask(int(taskID))
	if !failsucc{
		log.Printf("TestPassFailTask(): 2 %v")
		return false
	}
	task3, _, _:= GetTaskId(int(taskID))
	if task3.Status != "failed" {
		log.Printf("TestPassFailTask(): bad value on true fal%v", task3.Status)
		return false
	}
	return true


}

func TestDeleteTask() bool {
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
	// tx, err := DB.Beginx()
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

	// var userExists bool
	// err = tx.Get(&userExists, "SELECT EXISTS (SELECT 1 FROM UserTable WHERE UserID = ?)", newTask.UserID)
	// if err != nil {
	// 	fmt.Println("CreateTask(): breaky 2", err)
	// 	return false
	// }

	// if !userExists {
	// 	_, err = tx.Exec("INSERT INTO UserTable (UserID, Points, BossId) VALUES (?, 0, 2)", newTask.UserID)
	// 	if err != nil {
	// 		fmt.Println("CreateTask(): breaky 3", err)
	// 		return false
	// 	}
	// }

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

	taskResult, _, _ := GetTaskId(int(taskID))

	if taskResult.TaskName != "edited name" ||
		taskResult.Description != "edited description" ||
		taskResult.Status != "failed" ||
		taskResult.IsAllDay != true ||
		taskResult.Difficulty != "medium" {
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

func TestGetUserTaskTime() bool {

	starttime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(time.Hour)
	taskl, err := GetUserTaskDateTime(testUserId, starttime, endTime)
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
