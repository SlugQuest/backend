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
	TaskID         int
	UserID         string
	Category       string
	TaskName       string
	Description    string
	StartTime      time.Time
	EndTime        time.Time
	Status         string
	IsRecurring    bool
	IsAllDay       bool
	RecurringType  string
	Difficulty     string
	CronExpression string
}

type TaskPreview struct {
	TaskID      int
	UserID      string
	Category    string
	TaskName    string
	StartTime   time.Time
	EndTime     time.Time
	Status      string
	IsRecurring bool
	IsAllDay    bool
}

var DB *sqlx.DB

func LoadDumbData() error {
	// No recur patterns since we aren't using them yet
	for i := 1000; i < 1500; i++ {
		task := Task{TaskID: i, UserID: "1111", Category: "yo", TaskName: "some name" + strconv.Itoa(i), Description: "sumdesc" + strconv.Itoa(i), StartTime: time.Now(), EndTime: time.Now(), Status: "todo", IsRecurring: false, IsAllDay: false, CronExpression: "dummycron", Difficulty: "easy"}
		lol, _, err := CreateTask(task)
		if !lol || (err != nil) {
			return err
		}
	}
	return nil
}

func ConnectToDB(isunittest bool) error {
	if isunittest {
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
func passtask(Tid int)(error){
	stmt, err := DB.Preparex(`
	UPDATE TaskTable 
	SET Status = ?
	WHERE TaskID = ?
	`)
	if err != nil {
		stmt.Close()
		return err
	}
	_, erro := stmt.Exec("completed")
	stmt.Close()
	return erro;

}

func failtask(Tid int)(error){
	stmt, err := DB.Preparex(`
	UPDATE TaskTable 
	SET Status = ?
	WHERE TaskID = ?
	`)
	if err != nil {
		stmt.Close()
		return err
	}
	_, erro := stmt.Exec("failed")
	stmt.Close()
	return erro;

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

func CreateTask(task Task) (bool, int64, error) {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		fmt.Println("CreateTask(): breaky 1")
		return false, -1, err
	}
	defer tx.Rollback() //abort transaction if error

	//preparing statement to prevent SQL injection issues
	stmt, err := tx.Preparex("INSERT INTO TaskTable (UserID, Category, TaskName, Description, StartTime, EndTime, Status, IsRecurring, IsAllDay, Difficulty, CronExpression) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("CreateTask(): breaky 2")
		return false, -1, err
	}

	defer stmt.Close() // Defer the closing of SQL statement to ensure it closes once the function completes
	res, err := stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression)

	if err != nil {
		fmt.Println("CreateTask(): breaky 3 ", err)
		return false, -1, err
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("CreateTask(): breaky 4 ", err)
		return false, -1, err
	}

	points := CalculatePoints(task.Difficulty)
	_, err = tx.Exec("UPDATE UserTable SET Points = Points + ? WHERE UserID = ?", points, task.UserID)
	if err != nil {
		fmt.Println("CreateTask(): breaky 5 ", err)
		return false, -1, err
	}

	// if task.IsRecurring {
	// 	rStmnt, err := tx.Preparex("INSERT INTO RecurrencePatterns (TaskID, RecurringType, DayOfWeek, DayOfMonth) VALUES (?, ?, ?, ?)")
	// 	if err != nil {
	// 		fmt.Println("CreateTask(): breaky 4", err)
	// 		return false, -1, err
	// 	}
	// 	defer rStmnt.Close()

	// 	_, err = rStmnt.Exec(taskID, task.RecurringType, task.DayOfWeek, task.DayOfMonth)
	// 	if err != nil {
	// 		fmt.Println("CreateTask(): breaky 5", err)
	// 		return false, -1, err
	// 	}
	// }

	tx.Commit() //commit transaction to database

	return true, taskID, nil
}

func EditTask(task Task, id int) (bool, error) {

	tx, err := DB.Beginx()
	if err != nil {
		return false, err
	}

	var currentDifficulty string
	err = tx.Get(&currentDifficulty, "SELECT Difficulty FROM TaskTable WHERE TaskID = ?", id)
	if err != nil {
		tx.Rollback()
		fmt.Println("EditTask(): breaky 1", err)
		return false, err
	}

	stmt, err := tx.Preparex(`
		UPDATE TaskTable 
		SET UserID = ?, Category = ?, TaskName = ?, Description = ?, StartTime = ?, EndTime = ?, Status = ?, IsRecurring = ?, IsAllDay = ?, Difficulty = ?, CronExpression = ? 
		WHERE TaskID = ?
	`)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression, id)

	if err != nil {
		return false, err
	}

	oldPoints := CalculatePoints(currentDifficulty)
	newPoints := CalculatePoints(task.Difficulty)
	_, err = tx.Exec("UPDATE UserTable SET Points = Points - ? + ? WHERE UserID = ?", oldPoints, newPoints, task.UserID)
	if err != nil {
		fmt.Println("EditTask(): breaky 2", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteTask(id int) (bool, error) {
	tx, err := DB.Beginx()

	if err != nil {
		return false, err
	}

	// recurrenceTableExists, err := isTableExists("RecurrencePatterns")
	// if err != nil {
	// 	tx.Rollback()
	// 	fmt.Println("in here 1")
	// 	return false, err
	// }

	// if recurrenceTableExists {
	// 	stmt, err := tx.Preparex("DELETE FROM RecurrencePatterns WHERE TaskID = ?")
	// 	if err != nil {
	// 		tx.Rollback()
	// 		fmt.Println("in here 2", err)
	// 		return false, err
	// 	}
	// 	defer stmt.Close()

	// 	_, err = stmt.Exec(id)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		fmt.Println("in here 3", err)
	// 		return false, err
	// 	}
	// }

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
	rows, err := DB.Query("SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, Status, IsRecurring, IsAllDay FROM TaskTable;")
	utaskArr := []TaskPreview{}
	if err != nil {
		fmt.Println(err)
		rows.Close()
		return utaskArr, err
	}

	for rows.Next() {
		var taskprev TaskPreview
		erro := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay)
		if erro != nil {
			fmt.Println(erro)
			rows.Close()
		}
		utaskArr = append(utaskArr, taskprev)
	}
	rows.Close()
	return utaskArr, err
}
func GetUserTaskDateTime(Uid string, startq time.Time, endq time.Time) ([]TaskPreview, error) {
	prep, err := DB.Preparex("SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, Status, IsRecurring, IsAllDay FROM TaskTable t WHERE t.StartTime > ? AND t.EndTime < ?;")
	utaskArr := []TaskPreview{}
	if err != nil {
		fmt.Println(err)
		return utaskArr, err
	}
	rows, erro := prep.Query(startq, endq)
	if erro != nil {
		fmt.Println(err)
		rows.Close()
		return utaskArr, err
	}
	for rows.Next() {
		var taskprev TaskPreview
		erro := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay)
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
func GetTaskId(Tid int) (Task, error) {
	rows  := DB.QueryRow("SELECT * FROM TaskTable WHERE TaskID=?;", Tid)
	var taskit Task
	err := rows.Scan(&taskit.TaskID, &taskit.UserID, &taskit.Category, &taskit.TaskName, &taskit.Description, &taskit.StartTime, &taskit.EndTime, &taskit.Status, &taskit.IsRecurring, &taskit.IsAllDay, &taskit.Difficulty, &taskit.CronExpression)
	fmt.Println("finding")
	fmt.Println("done finding")
	fmt.Println(taskit.Status)
	return taskit, err
}
