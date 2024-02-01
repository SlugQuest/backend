package crud

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	TaskID        int
	UserID        string
	Category      string
	TaskName      string
	Description   string
	StartTime     time.Time
	EndTime       time.Time
	IsCompleted   bool
	IsRecurring   bool
	IsAllDay      bool
	RecurringType string
	DayOfWeek     int
	DayOfMonth    int
}

type TaskPreview struct {
	TaskID      int
	UserID      string
	Category    string
	TaskName    string
	StartTime   time.Time
	EndTime     time.Time
	IsCompleted bool
	IsRecurring bool
	IsAllDay    bool
}

var DB *sqlx.DB

func LoadDumbData() error {
	// No recur patterns since we aren't using them yet
	for i := 1000; i < 1500; i++ {
		task := Task{TaskID: i, UserID: "1111", Category: "asdf", TaskName: "some name" + strconv.Itoa(i), Description: "sumdesc" + strconv.Itoa(i), StartTime: time.Now(), EndTime: time.Now(), IsCompleted: false, IsRecurring: false, IsAllDay: false}
		lol, err, _ := CreateTask(task)
		if !lol || (err != nil) {
			return err
		}
	}
	return nil

}

func ConnectToDB() error {
	// Read schema from file
	schemaCreate, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	fmt.Println(string(schemaCreate))

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

	// Execute the schema creation SQL
	_, err = db.Exec(string(schemaCreate))
	if err != nil {
		fmt.Println("Error executing schema creation SQL:", err)
		return err
	}

	DB = db
	return nil
}

func isTableExists(tableName string) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
	err := DB.Get(&count, query)
	return count > 0, err
}

func CreateTask(task Task) (bool, error, int64) {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		fmt.Println("breaky 1 ")
		return false, err, -1
	}
	defer tx.Rollback() //abort transaction if error

	//preparing statement to prevent SQL injection issues
	stmt, err := tx.Preparex("INSERT INTO TaskTable (UserID, Category, TaskName, Description, StartTime, EndTime, IsCompleted, IsRecurring, IsAllDay) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("breaky 2")
		return false, err, -1
	}

	defer stmt.Close() //defer the closing of SQL statement to ensure it Closes once the function completes
	res, err := stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.IsCompleted, task.IsRecurring, task.IsAllDay)

	if err != nil {
		fmt.Println("breaky 3 ", err)
		return false, err, -1
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("breaky 4 ", err)
		return false, err, -1
	}

	if task.IsRecurring {
		rStmnt, err := tx.Preparex("INSERT INTO RecurrencePatterns (TaskID, RecurringType, DayOfWeek, DayOfMonth) VALUES (?, ?, ?, ?)")
		if err != nil {
			fmt.Println("breaky 4", err)
			return false, err, -1
		}
		defer rStmnt.Close()

		_, err = rStmnt.Exec(taskID, task.RecurringType, task.DayOfWeek, task.DayOfMonth)
		if err != nil {
			fmt.Println("breaky 5", err)
			return false, err, -1
		}
	}

	tx.Commit() //commit transaction to database

	return true, nil, taskID
}

func EditTask(task Task, id int) (bool, error) {

	tx, err := DB.Beginx()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Preparex(`
		UPDATE TaskTable 
		SET UserID = ?, Category = ?, TaskName = ?, Description = ?, StartTime = ?, EndTime = ?, IsCompleted = ?, IsRecurring = ?, IsAllDay = ? 
		WHERE TaskID = ?
	`)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.IsCompleted, task.IsRecurring, task.IsAllDay, id)

	if err != nil {
		return false, err
	}

	// if the edit updated the recurrence of a task
	if task.IsRecurring {
		_, err = tx.Exec(
			`UPDATE RecurrencePatterns 
			SET RecurringType = ?, DayOfWeek = ?, DayOfMonth = ? 
			WHERE TaskID = ?
		`, task.RecurringType, task.DayOfWeek, task.DayOfMonth, id)

		if err != nil {
			tx.Rollback()
			return false, err
		}
	} else { // if isRecurring was set to false
		_, err = tx.Exec("DELETE FROM RecurrencePatterns WHERE TaskID = ?", id)
		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	tx.Commit()

	return true, nil
}

func DeleteTask(id int) (bool, error) {
	tx, err := DB.Beginx()

	if err != nil {
		return false, err
	}

	recurrenceTableExists, err := isTableExists("RecurrencePatterns")
	if err != nil {
		tx.Rollback()
		fmt.Println("in here 1")
		return false, err
	}

	if recurrenceTableExists {
		stmt, err := tx.Preparex("DELETE FROM RecurrencePatterns WHERE TaskID = ?")
		if err != nil {
			tx.Rollback()
			fmt.Println("in here 2", err)
			return false, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			fmt.Println("in here 3", err)
			return false, err
		}
	}

	stmt2, err := tx.Preparex("DELETE FROM TaskTable WHERE TaskID = ?")

	if err != nil {
		tx.Rollback()
		fmt.Println("in here 4", err)
		return false, err
	}

	defer stmt2.Close()

	_, err = stmt2.Exec(id)

	if err != nil {
		tx.Rollback()
		fmt.Println("in here 5", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}

// Need hardcode Uid for testing until we have auth0
func GetUserTask(Uid int) ([]TaskPreview, error) {
	rows, err := DB.Query("SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, IsCompleted, IsRecurring, IsAllDay FROM TaskTable;")
	utaskArr := []TaskPreview{}
	if err != nil {
		fmt.Println(err)
		rows.Close()
		return utaskArr, err
	}
	fmt.Println(Uid)
	for rows.Next() {
		var taskprev TaskPreview
		erro := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.IsCompleted, &taskprev.IsRecurring, &taskprev.IsAllDay)
		if erro != nil {
			fmt.Println(erro)
			rows.Close()
		}
		utaskArr = append(utaskArr, taskprev)
	}
	rows.Close()
	return utaskArr, err
}

func GetTaskId(Tid int) (Task, error, bool) {
	rows, err := DB.Query("SELECT * FROM TaskTable WHERE TaskID=?;", Tid)
	var taskit Task
	if err != nil {
		fmt.Println(err)
		return taskit, err, false
	}
	counter := 0
	for rows.Next() {
		counter += 1
		fmt.Println(counter)
		rows.Scan(&taskit.TaskID, &taskit.UserID, &taskit.Category, &taskit.TaskName, &taskit.Description, &taskit.StartTime, &taskit.EndTime, &taskit.IsCompleted, &taskit.IsRecurring, &taskit.IsAllDay)
		fmt.Println("finding")
	}
	rows.Close()
	fmt.Println("done finding")
	print(counter)
	return taskit, err, counter == 1
}
