package crud

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Shared within the package
var DB *sqlx.DB

// ID for the default team, for individual tasks
var NoTeamID int = -1

func LoadDumbData() error {
	// sql.NullInt64{Valid: false} -> Null for TeamID
	for i := 1000; i < 1500; i++ {
		task := Task{TaskID: i, UserID: "test_user_id", Category: "test_category", TaskName: "some name" + strconv.Itoa(i),
			Description: "sumdesc" + strconv.Itoa(i), StartTime: time.Now(), EndTime: time.Now(), Status: "todo", IsRecurring: false,
			IsAllDay: false, CronExpression: "dummycron", Difficulty: "easy", TeamID: NoTeamID}
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

	bossAdded, err := AddBoss(Boss{BossID: 1, Name: "testboss_name", Health: 30, Image: "../images/clown.jpeg"})
	if !bossAdded || err != nil {
		return err
	}

	return nil
}

func ConnectToDB(isUnitTest bool) error {
	if isUnitTest {
		// Read schema from file
		schemaCreate, err := os.ReadFile("../schema.sql")
		if err != nil {
			log.Printf("ConnectToDB (unit test): Error reading schema file: %v", err)
			return err
		}

		// Connect to an in-memory SQLite database
		db, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			log.Printf("ConnectToDB (unit test): Error opening in-memory database: %v", err)
			return err
		}

		// Force a connection and test that it worked
		err = db.Ping()
		if err != nil {
			log.Printf("ConnectToDB (unit test): Error pinging database: %v", err)
			return err
		}

		// Execute the schema creation SQL
		_, err = db.Exec(string(schemaCreate))
		if err != nil {
			log.Printf("ConnectToDB (unit test): Error executing schema creation SQL: %v", err)
			return err
		}

		DB = db
	} else {
		// Connect to the real database
		db, err := sqlx.Open("sqlite3", "slugquest.db")
		if err != nil {
			log.Printf("ConnectToDB: Error connecting to database: %v", err)
			return err
		}

		// Force a connection and test that it worked
		err = db.Ping()
		if err != nil {
			log.Printf("ConnectToDB: Error pinging database: %v", err)
			return err
		}

		DB = db
	}

	log.Println("Successfully connected to DB!")
	return nil
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

// for recurrence work
func CreateRecurringLogEntry(taskID int64, status string, timestamp time.Time) (bool, int64, error) {
	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("CreateRecurringLog(): breaky 1: %v\n", err)
		return false, -1, err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex("INSERT INTO RecurringLog (TaskID, Status, timestamp) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("CreateRecurringLog(): breaky 2: %v\n", err)
		return false, -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(taskID, status, timestamp)
	if err != nil {
		log.Printf("CreateRecurringLog(): breaky 3: %v\n", err)
		return false, -1, err
	}

	logID, err := res.LastInsertId()
	if err != nil {
		log.Printf("CreateRecurringLog(): breaky 4: %v\n", err)
		return false, -1, err
	}

	tx.Commit()

	return true, logID, nil
}

func GetRecurringTasks() ([]Task, error) {
	var recurringTasks []Task

	query := `SELECT * FROM TaskTable WHERE IsRecurring = true`

	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("GetRecurringTasks(): Error querying recurring tasks: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.TaskID,
			&task.UserID,
			&task.Category,
			&task.TaskName,
			&task.Description,
			&task.StartTime,
			&task.EndTime,
			&task.Status,
			&task.IsRecurring,
			&task.IsAllDay,
			&task.Difficulty,
			&task.CronExpression,
			&task.TeamID,
		)
		if err != nil {
			log.Printf("GetRecurringTasks(): Error scanning row: %v\n", err)
			return nil, err
		}
		recurringTasks = append(recurringTasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetRecurringTasks(): Error iterating over rows: %v\n", err)
		return nil, err
	}

	if len(recurringTasks) == 0 {
		log.Println("No recurring tasks found.")
	}

	return recurringTasks, nil
}

func PopRecurringTasksMonth() error {
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	recurringTasks, err := GetRecurringTasks()
	if err != nil {
		return err
	}

	for _, task := range recurringTasks {
		cronExpression := task.CronExpression
		log.Printf("Parsing cron expression: %s\n", cronExpression)

		nextTimes := cronexpr.MustParse(task.CronExpression).NextN(time.Now(), 31)
		//assuming there can only be one recurrence a day, so at most 31 recurrences in a month

		for _, nextTime := range nextTimes {
			// Check if the next occurrence is in the current month
			if nextTime.Month() == currentMonth && nextTime.Year() == currentYear {
				_, _, err = CreateRecurringLogEntry(int64(task.TaskID), "todo", nextTime)
				if err != nil {
					log.Printf("PopRecurringTasksMonth(): %v", err)
					return err
				}
			}
		}
	}
	return nil
}

func CountRecurringLogEntries() (int, error) {
	var count int

	query := "SELECT COUNT(*) FROM RecurringLog"

	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("CountRecurringLogEntries(): Error executing query: %v\n", err)
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Printf("CountRecurringLogEntries(): Error scanning row: %v\n", err)
			return 0, err
		}
	} else {
		log.Println("CountRecurringLogEntries(): No rows returned.")
		return 0, nil
	}

	log.Printf("Number of recurring log entries: %d\n", count)
	return count, nil
}

func AddDefaultTeam() (bool, error) {
	_, err := DB.Exec("INSERT OR IGNORE INTO TeamTable VALUES (?, ?)", NoTeamID, "NoTeam")
	if err != nil {
		log.Printf("AddDefaultTeam(): error adding default team to DB: %v", err)
		return false, err
	}
	return true, nil
}
