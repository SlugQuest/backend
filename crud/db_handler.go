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

type User struct {
	UserID   string // Not known to user, do not expose
	Username string // Set by user, can be exposed
	Picture  string // A0 stores their profile pics as URLs
	Points   int
	BossId   int
}

var DB *sqlx.DB

func LoadDumbData() error {
	// No recur patterns since we aren't using them yet
	for i := 1000; i < 1500; i++ {
		task := Task{TaskID: i, UserID: "1111", Category: "asdf", TaskName: "some name" + strconv.Itoa(i), Description: "sumdesc" + strconv.Itoa(i), StartTime: time.Now(), EndTime: time.Now(), IsCompleted: false, IsRecurring: false, IsAllDay: false}
		lol, _, err := CreateTask(task)
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

func CreateTask(task Task) (bool, int64, error) {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		fmt.Println("CreateTask(): breaky 1")
		return false, -1, err
	}
	defer tx.Rollback() //abort transaction if error

	//preparing statement to prevent SQL injection issues
	stmt, err := tx.Preparex("INSERT INTO TaskTable (UserID, Category, TaskName, Description, StartTime, EndTime, IsCompleted, IsRecurring, IsAllDay) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("CreateTask(): breaky 2")
		return false, -1, err
	}

	defer stmt.Close() //defer the closing of SQL statement to ensure it Closes once the function completes
	res, err := stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.IsCompleted, task.IsRecurring, task.IsAllDay)

	if err != nil {
		fmt.Println("CreateTask(): breaky 3 ", err)
		return false, -1, err
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("CreateTask(): breaky 4 ", err)
		return false, -1, err
	}

	if task.IsRecurring {
		rStmnt, err := tx.Preparex("INSERT INTO RecurrencePatterns (TaskID, RecurringType, DayOfWeek, DayOfMonth) VALUES (?, ?, ?, ?)")
		if err != nil {
			fmt.Println("CreateTask(): breaky 4", err)
			return false, -1, err
		}
		defer rStmnt.Close()

		_, err = rStmnt.Exec(taskID, task.RecurringType, task.DayOfWeek, task.DayOfMonth)
		if err != nil {
			fmt.Println("CreateTask(): breaky 5", err)
			return false, -1, err
		}
	}

	tx.Commit() //commit transaction to database

	return true, taskID, nil
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

// Uid is provided in a router context (session cookies)
func GetUserTask(Uid string) ([]TaskPreview, error) {
	rows, err := DB.Query("SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, IsCompleted, IsRecurring, IsAllDay FROM TaskTable;")
	utaskArr := []TaskPreview{}
	if err != nil {
		fmt.Println(err)
		rows.Close()
		return utaskArr, err
	}

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

// Find task by TaskID
func GetTaskId(Tid int) (Task, bool, error) {
	rows, err := DB.Query("SELECT * FROM TaskTable WHERE TaskID=?;", Tid)
	var taskit Task
	if err != nil {
		fmt.Println(err)
		return taskit, false, err
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
	return taskit, counter == 1, err
}

// Find user by UserID
func GetUser(Uid string) (User, bool, error) {
	rows, err := DB.Query("SELECT * FROM UserTable WHERE UserID=?;", Uid)
	var user User
	if err != nil {
		fmt.Println(err)
		return user, false, err
	}

	counter := 0
	for rows.Next() {
		counter += 1
		rows.Scan(&user.UserID, &user.Points, &user.BossId)
	}
	rows.Close()

	return user, counter == 1, err
}

// Add user into DB
func AddUser(u User) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		fmt.Println("AddUser(): breaky 1")
		return false, err
	}
	defer tx.Rollback() // aborrt transaction if error

	stmt, err := tx.Preparex("INSERT INTO UserTable (UserID, Points, Bossid) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("AddUser(): breaky 2")
		return false, err
	}

	defer stmt.Close() //defer the closing of SQL statement to ensure it Closes once the function completes
	_, err = stmt.Exec(u.UserID, u.Points, u.BossId)
	if err != nil {
		fmt.Println("AddUser(): breaky 3")
		return false, err
	}

	tx.Commit() //commit transaction to database

	return true, nil
}

// Edit a user by supplying new values
func EditUser(u User, uid string) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Preparex(`
		UPDATE UserTable 
		SET UserID = ?, Points = ?, Bossid = ?
		WHERE UserID = ?
	`)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.UserID, u.Points, u.BossId, uid)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
