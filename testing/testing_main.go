package testing

// When a new backend function is made, add a test function for it that returns a bool, and then put that func in testmain
import (
	"log"
	"os"
	"testing"

	"slugquest.com/backend/crud"
)

// var dummyUserID string = "1111"
// var testUserID string = "2222" // testing user functions

// func RunAllTests() bool {
// 	ConnectToDB(true)
// 	dummy_err := LoadDumbData()
// 	if dummy_err != nil {
// 		log.Fatalf("error loaduing dumb data: %v", dummy_err)
// 	}
// 	return TestGetCurrBossHealth() && TestGetUserTask() && TestGetCategory() && TestDeleteTask() && TestPassFailTask() && TestEditTask() && TestGetTaskId() && TestAddUser() && TestEditUser() && TestDeleteUser()
// }

func TestMain(m *testing.M) {
	// Setup
	conn_err := crud.ConnectToDB(true)
	if conn_err != nil || crud.DB == nil {
		log.Fatalf("Error setting up DB for unit tests: %v", conn_err)
	}

	dummy_err := crud.LoadDumbData()
	if dummy_err != nil {
		log.Fatalf("error loaduing dumb data: %v", dummy_err)
	}

	// Run all tests in this package
	result_code := m.Run()

	// Teardown
	crud.DB.Close()

	os.Exit(result_code)
}
