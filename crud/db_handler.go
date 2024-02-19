package crud

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Shared within the package
var DB *sqlx.DB

func LoadDumbData() error {
	// No recur patterns since we aren't using them yet
	for i := 1000; i < 1500; i++ {
		task := Task{TaskID: i, UserID: "test_user_id", Category: "yo", TaskName: "some name" + strconv.Itoa(i), Description: "sumdesc" + strconv.Itoa(i), StartTime: time.Now(), EndTime: time.Now(), Status: "todo", IsRecurring: false, IsAllDay: false, CronExpression: "dummycron", Difficulty: "easy"}
		lol, _, err := CreateTask(task)
		if !lol || (err != nil) {
			return err
		}
	}
	for i := 1000; i < 1500; i++ {
		cat := Category{CatID: i, UserID: "test_user_id", Name: "lolcat", Color: 255}
		lol2, _, err2 := CreateCategory(cat)
		if !lol2 || (err2 != nil) {
			return err2
		}
	}
	return nil
}

func ConnectToDB(isunittest bool) error {
	if isunittest {
		// Read schema from file
		schemaCreate, err := os.ReadFile("../schema.sql")
		if err != nil {
			return err
		}

		// Connect to an in-memory SQLite database
		db, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			return err
		}

		// Force a connection and test that it worked
		err = db.Ping()
		if err != nil {
			fmt.Println("breaky")
			return err
		} else {
			fmt.Println("not breaky")
		}

		//Execute the schema creation SQL
		_, err = db.Exec(string(schemaCreate))
		if err != nil {
			fmt.Println("Error executing schema creation SQL:", err)
			return err
		}

		DB = db
	} else {

		// Connect to the real database
		db, err := sqlx.Open("sqlite3", "slugquest.db")
		if err != nil {
			return err
		}

		// Force a connection and test that it worked
		err = db.Ping()
		if err != nil {
			fmt.Println("breaky")
			return err
		} else {
			fmt.Println("not breaky")
		}

		DB = db
	}
	return nil
}

func isTableExists(tableName string) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
	err := DB.Get(&count, query)
	return count > 0, err
}

func CalculatePoints(difficulty string) int {
	switch difficulty {
	case "easy":
		return 1
	case "medium":
		return 2
	case "hard":
		return 3
	default:
		return 0
	}
}

// for recurrence work in the future
func CreateRecurringLogEntry(taskID int, isCurrent bool, status string) (bool, int64, error) {
	tx, err := DB.Beginx()
	if err != nil {
		fmt.Printf("CreateRecurringLog(): breaky 1: %v\n", err)
		return false, -1, err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex("INSERT INTO RecurringLog (TaskID, isCurrent, Status) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Printf("CreateRecurringLog(): breaky 2: %v\n", err)
		return false, -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(taskID, isCurrent, status)
	if err != nil {
		fmt.Printf("CreateRecurringLog(): breaky 3: %v\n", err)
		return false, -1, err
	}

	logID, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("CreateRecurringLog(): breaky 4: %v\n", err)
		return false, -1, err
	}

	tx.Commit()

	return true, logID, nil
}
