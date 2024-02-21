package crud

import (
	"fmt"
	"log"
	"time"

	_ "github.com/gorhill/cronexpr"
)

// Find task by TaskID
func GetTaskId(tid int) (Task, bool, error) {
	var taskit Task

	prep, err := DB.Preparex("SELECT * FROM TaskTable WHERE TaskID=?;")
	if err != nil {
		log.Printf("GetTaskId() #1: %v", err)
		return taskit, false, err
	}
	defer prep.Close()

	rows, err := prep.Query(tid)
	if err != nil {
		log.Printf("GetTaskId() #2: %v", err)
		rows.Close()
		return taskit, false, err
	}

	counter := 0
	for rows.Next() {
		counter += 1
		err := rows.Scan(&taskit.TaskID, &taskit.UserID, &taskit.Category, &taskit.TaskName, &taskit.Description, &taskit.StartTime, &taskit.EndTime, &taskit.Status, &taskit.IsRecurring, &taskit.IsAllDay, &taskit.Difficulty, &taskit.CronExpression)
		if err != nil {
			log.Printf("GetTaskId() #3: %v", err)
			rows.Close()
		}
	}

	prep.Close()
	rows.Close()
	return taskit, counter == 1, err
}

// Uid is provided in a router context (session cookies)
func GetUserTask(uid string) ([]TaskPreview, error) {
	utaskArr := []TaskPreview{}

	prep, err := DB.Preparex(`SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, Status, IsRecurring, IsAllDay FROM TaskTable
		WHERE UserID = ?;`)
	if err != nil {
		log.Printf("GetUserTask() #1: %v", err)
		return utaskArr, err
	}

	rows, err := prep.Query(uid)
	if err != nil {
		log.Printf("GetUserTask() #2: %v", err)
		rows.Close()
		prep.Close()
		return utaskArr, err
	}

	for rows.Next() {
		var taskprev TaskPreview
		err := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay)
		if err != nil {
			log.Printf("GetUserTask() #3: %v", err)
			rows.Close()
		}
		utaskArr = append(utaskArr, taskprev)
	}
	prep.Close()
	rows.Close()
	return utaskArr, err
}

func GetUserTaskDateTime(uid string, startq time.Time, endq time.Time) ([]TaskPreview, error) {
	utaskArr := []TaskPreview{}

	prep, err := DB.Preparex(`SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, Status, IsRecurring, IsAllDay FROM TaskTable t 
		WHERE UserID = ? AND t.StartTime > ? AND t.EndTime < ?;`)
	if err != nil {
		log.Printf("GetUserTaskDateTime() #1: %v", err)
		return utaskArr, err
	}
	defer prep.Close()

	rows, err := prep.Query(uid, startq, endq)
	if err != nil {
		log.Printf("GetUserTaskDateTime() #2: %v", err)
		rows.Close()
		return utaskArr, err
	}

	for rows.Next() {
		var taskprev TaskPreview
		err := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay)
		if err != nil {
			fmt.Println(err)
			rows.Close()
			return utaskArr, err
		}
		utaskArr = append(utaskArr, taskprev)
	}
	prep.Close()
	rows.Close()
	// p2, err := DB.Preparex("SELECT * FROM RecurringLog WHERE UserID = ? AND t.StartTime > ? AND t.EndTime < ?;")
	// if err != nil {
	// 	log.Printf("GetUserTaskDateTime() #2: %v", err)
	// 	return utaskArr, err
	// }
	// rowrec, err := prep.Query

	return utaskArr, err
}

func CreateTask(task Task) (bool, int64, error) {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		fmt.Println("CreateTask(): breaky 1")
		return false, -1, err
	}
	defer tx.Rollback() // Abort transaction if any error occurs

	//preparing statement to prevent SQL injection issues
	stmt, err := tx.Preparex("INSERT INTO TaskTable (UserID, Category, TaskName, Description, StartTime, EndTime, Status, IsRecurring, IsAllDay, Difficulty, CronExpression) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("CreateTask(): breaky 2", err)
		return false, -1, err
	}

	defer stmt.Close() // Defer the closing of SQL statement to ensure it closes once the function completes
	fmt.Println(task.UserID)
	res, err := stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression)

	if err != nil {
		// fmt.Println(task)
		fmt.Println("CreateTask(): breaky 3 ", err)
		return false, -1, err
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("CreateTask(): breaky 4 ", err)
		return false, -1, err
	}

	// if task.IsRecurring {
	// 	// rStmnt, err := tx.Preparex("INSERT INTO RecurrencePatterns (TaskID, RecurringType, DayOfWeek, DayOfMonth) VALUES (?, ?, ?, ?)")
	// 	// if err != nil {
	// 	// 	fmt.Println("CreateTask(): breaky 4", err)
	// 	// 	return false, -1, err
	// 	// }
	// 	nexTime := cronexpr.MustParse("0 0 1 * * ?").NextN(task.StartTime, 10)
	// 	// fmt.Println(nexTime)
	// 	rStmnt, err := tx.Preparex("INSERT INTO RecurringLog (TaskID, isCurrent, Status, CreatedAt) VALUES (?, ?, ?, ?)")
	// 	if err != nil {
	// 		fmt.Println("CreateTask(): breaky 4", err)
	// 		return false, -1, err
	// 	}
	// 	for i := 0; i < 5; i++ {
	// 		_, err := rStmnt.Exec(taskID, false, task.Status,nexTime[i] )

	// 		if err != nil {
	// 			// fmt.Println(task)
	// 			fmt.Println("CreateTask(): breaky 7 ", err)
	// 			return false, -1, err
	// 		}
	// 	}

	// }
	// 	defer rStmnt.Close()

	// 	_, err = rStmnt.Exec(taskID, task.RecurringType, task.DayOfWeek, task.DayOfMonth)
	// 	if err != nil {
	// 		fmt.Println("CreateTask(): breaky 5", err)
	// 		return false, -1, err
	// 	}
	// }

	tx.Commit() //commit transaction to database
	// fmt.Println("WE ADDED A TASK")
	return true, taskID, nil
}

func EditTask(task Task, tid int, uid string) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("EditTask() #1: %v", err)
		return false, err
	}

	stmt, err := tx.Preparex(`
		UPDATE TaskTable 
		SET Category = ?, TaskName = ?, Description = ?, StartTime = ?, EndTime = ?, Status = ?, IsRecurring = ?, IsAllDay = ?, Difficulty = ?, CronExpression = ? 
		WHERE TaskID = ? AND UserID = ?
	`)

	if err != nil {
		log.Printf("EditTask() #2: %v", err)
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression,
		tid, uid)
	if err != nil {
		log.Printf("EditTask() #3: %v", err)
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("EditTask() #4: %v", err)
		return false, err
	}

	return true, nil
}

func DeleteTask(tid int, uid string) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("DeleteTask() #1: %v", err)
		return false, err
	}
	defer tx.Rollback() // Abort transaction if any error occurs

	// recurrenceTableExists, err := isTableExists("RecurrencePatterns")
	// if err != nil {
	// 	fmt.Println("in here 1")
	// 	return false, err
	// }

	delTT, err := tx.Preparex("DELETE FROM TaskTable WHERE TaskID = ? AND UserID = ?")
	if err != nil {
		log.Printf("DeleteTask() #2: %v", err)
		return false, err
	}
	defer delTT.Close()

	_, err = delTT.Exec(tid)
	if err != nil {
		log.Printf("DeleteTask() #3: %v", err)
		return false, err
	}

	delRL, err := tx.Preparex("DELETE FROM RecurringLog WHERE TaskID = ?")
	if err != nil {
		log.Printf("DeleteTask() #4: can't preparing statement for RecurringLog deletion: %v", err)
		return false, err
	}
	defer delRL.Close()

	_, err = delRL.Exec(tid)
	if err != nil {
		log.Printf("DeleteTask() #5: Error deleting from RecurringLog: %v", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}

func Passtask(Tid int) (bool, error) {

	task, ok, err := GetTaskId(Tid)
	if err != nil {
		fmt.Printf("Passtask(): breaky 2 %v\n", err)
		return false, err
	}

	if !ok {
		fmt.Println("Passtask(): Task not found")
		return false, fmt.Errorf("task not found")
	}

	if task.IsRecurring {
		_, err := DB.Exec(`
			UPDATE RecurringLog 
			SET Status = ?
			WHERE TaskID = ? AND isCurrent = true
		`, "completed", Tid)

		if err != nil {
			fmt.Printf("Passtask(): breaky 0 %v\n", err)
			return false, err
		}

	} else {
		tx, err := DB.Beginx() // start transaction
		if err != nil {
			fmt.Printf("Passtask(): breaky 1 %v\n", err)
			return false, err
		}
		defer tx.Rollback() // Abort transaction if any error occurs
		stmt, err := tx.Preparex(`
			UPDATE TaskTable 
			SET Status = ?
			WHERE TaskID = ?
		`)

		if err != nil {
			fmt.Printf("Passtask(): breaky 2 %v\n", err)
			return false, err
		}

		_, err = stmt.Exec("completed", Tid)
		if err != nil {
			fmt.Printf("Passtask(): breaky 3 %v\n", err)
			return false, err
		}

		tx.Commit()

	}

	// tx, err = DB.Beginx() // start transaction
	// if err != nil {
	// 	fmt.Printf("Passtask(): breaky %v\n", err)
	// 	return false, err
	// }

	points := CalculatePoints(task.Difficulty)
	_, err = DB.Exec("UPDATE UserTable SET Points = Points + ? WHERE UserID = ?", points, task.UserID)
	if err != nil {
		fmt.Printf("Passtask(): breaky 5 %v\n", err)
		return false, err
	}

	currBossHealth, err := GetCurrBossHealth(task.UserID)
	if err != nil {
		fmt.Printf("Passtask(): breaky %v\n", err)
		return false, err
	}

	// Check if the current boss health is zero
	if currBossHealth <= 0 {
		// Switch to the next boss ID (currBossId + 1)
		_, err := DB.Exec("UPDATE UserTable SET BossId = BossId + 1 WHERE UserID = ?", task.UserID)
		if err != nil {
			fmt.Printf("Passtask(): breaky 6 %v\n", err)
			return false, err
		}

		// Reset user points to 0
		_, err = DB.Exec("UPDATE UserTable SET Points = ? WHERE UserID = ?", 0, task.UserID)
		if err != nil {
			fmt.Printf("Passtask(): breaky 7 %v\n", err)
			return false, err
		}
	}

	return true, nil
}

func Failtask(Tid int) bool {
	task, ok, err := GetTaskId(Tid)
	if err != nil {
		fmt.Printf("Failtask(): breaky %v\n", err)
		return false
	}

	if !ok {
		fmt.Println("Failtask(): Task not found")
		return false
	}

	if task.IsRecurring {
		_, err := DB.Exec(`
			UPDATE RecurringLog 
			SET Status = ?
			WHERE TaskID = ? AND isCurrent = true
		`, "failed", Tid)

		if err != nil {
			fmt.Printf("Failtask(): breaky 0 %v\n", err)
			return false
		}
	} else {
		tx, err := DB.Beginx() //start transaction
		if err != nil {
			return false
		}
		defer tx.Rollback() // Abort transaction if any error occurs

		stmt, err := tx.Preparex(`
		UPDATE TaskTable 
		SET Status = ?
		WHERE TaskID = ?
		`)
		if err != nil {
			return false
		}
		swag, erro := stmt.Exec("failed", Tid)
		stmt.Close()
		if erro != nil {
			print(erro.Error())
			print("FailtTask(): breaky 1 ")
			fmt.Println(erro)
			fmt.Println(swag)
			return false
		}

		tx.Commit()
	}

	return true

}
