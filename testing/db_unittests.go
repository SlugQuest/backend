package testing

// When a new backend function is made, add a test function for it that returns a bool, and then put that func in testmain
import (
	"log"
	"time"

	. "slugquest.com/backend/crud"
)

var dummyUserID string = "1111" // testing task functions
var testUserID string = "2222"  // testing user functions

func RunAllTests() bool {
	return TestGetUserTask() && TestDeleteTask() && TestEditTask() && TestGetTaskId() && TestAddUser() && TestDeleteUser()
}

func TestDeleteTask() bool {
	newTask := Task{
		UserID:      dummyUserID,
		Category:    "example",
		TaskName:    "New Task",
		Description: "Description of the new task",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		IsCompleted: false,
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
		UserID:      dummyUserID,
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

	success, taskID, err := CreateTask(newTask)
	if err != nil || !success {
		log.Printf("TestEditTask(): error creating task: %v", err)
		return false
	}

	editedTask := Task{
		TaskID:        int(taskID),
		UserID:        dummyUserID,
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
		log.Printf("TestEditTask(): error editing task: %v", editErr)
		return false
	}

	taskl, _, _ := GetTaskId(int(taskID))
	if taskl.TaskName != "edited name" || !taskl.IsCompleted {
		log.Println("TestEditTask(): edit verfication failed")
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
		log.Printf("TestDeleteUser(): couldn't add user")
		return false
	}

	_, found, _ := GetUser(newUser.UserID)
	if !found {
		log.Println("TestDeleteUser(): add failed")
		return false
	}

	return true
}

func TestEditUser() bool {
	// Original is one inserted in TestAddUser()
	editedUser := User{
		UserID:   testUserID,
		Username: "edited name",
		Picture:  "different.jpg",
		Points:   5,
		BossId:   10,
	}

	editSuccess, editErr := EditUser(editedUser, editedUser.UserID)
	if editErr != nil || !editSuccess {
		log.Printf("TestEditUser(): error editing user: %v", editErr)
		return false
	}

	checkE, _, _ := GetUser(editedUser.UserID)
	if checkE.Username != editedUser.Username || checkE.Picture != editedUser.Picture || !(checkE.Points == 5) || !(checkE.BossId == 10) {
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
