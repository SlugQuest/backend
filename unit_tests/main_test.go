package unit_tests

// When a new backend function is made, add a test function for it that returns a bool, and then put that func in testmain
import (
	"log"
	"os"
	"testing"

	"slugquest.com/backend/crud"
)

var testUser = crud.User{
	UserID:   "test_user_id",
	Username: "sluggo1",
	Picture:  "lol.jpg",
	Points:   1,
	BossId:   1,
}

func TestMain(m *testing.M) {
	setupDBForUnitTests()

	// Run all tests in this package
	result_code := m.Run()

	teardownAfterTests()

	os.Exit(result_code)
}

func setupDBForUnitTests() {
	conn_err := crud.ConnectToDB(true)
	if conn_err != nil || crud.DB == nil {
		log.Fatalf("Error setting up DB for unit tests: %v", conn_err)
	}

	dummy_err := crud.LoadDumbData()
	if dummy_err != nil {
		log.Fatalf("error loaduing dumb data: %v", dummy_err)
	}
	crud.AddUser(testUser)
}

func teardownAfterTests() {
	crud.DeleteUser(userForUserTable.UserID)
	crud.DB.Close()
}
