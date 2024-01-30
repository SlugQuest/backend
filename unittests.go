package main

// When a new backend function is made, add a test function for it that returns a bool, and then put that func in testmain
import (
	"fmt"
)

func testmain() bool {
	return TestCreateTask() && TestEditTask() && TestGetUserTask() && TestGetTaskId()
}


func TestCreateTask() bool {
	return true
}
func TestEditTask() bool {
	return true
}
func TestGetUserTask() bool {
	taskl, err := GetUserTask(1111)
	if err != nil{
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