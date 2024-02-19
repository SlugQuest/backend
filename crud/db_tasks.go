package crud

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

// Find task by TaskID
func GetTaskId(Tid int) (Task, bool, error) {
	rows, err := DB.Queryx("SELECT * FROM TaskTable WHERE TaskID=?;", Tid)
	var taskit Task
	if err != nil {
		fmt.Println(err)
		return taskit, false, err
	}
	counter := 0
	for rows.Next() {
		counter += 1
		fmt.Println(counter)
		rows.Scan(&taskit.TaskID, &taskit.UserID, &taskit.Category, &taskit.TaskName, &taskit.Description, &taskit.StartTime, &taskit.EndTime, &taskit.Status, &taskit.IsRecurring, &taskit.IsAllDay, &taskit.Difficulty, &taskit.CronExpression)
		fmt.Println("finding")
	}
	rows.Close()
	fmt.Println("done finding")
	fmt.Println(counter)
	fmt.Println(taskit.Status)
	return taskit, counter == 1, err
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
		prep.Close()
		return utaskArr, err
	}
	rows, erro := prep.Query(startq, endq)
	if erro != nil {
		fmt.Println(err)
		rows.Close()
		prep.Close()
		return utaskArr, err
	}
	for rows.Next() {
		var taskprev TaskPreview
		erro := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay)
		if erro != nil {
			fmt.Println(erro)
			rows.Close()
			prep.Close()
			return utaskArr, erro
		}
		utaskArr = append(utaskArr, taskprev)
	}
	prep.Close()
	rows.Close()
	return utaskArr, err
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
		fmt.Println("CreateTask(): breaky 2", err)
		return false, -1, err
	}

	defer stmt.Close() // Defer the closing of SQL statement to ensure it closes once the function completes
	fmt.Println(task)
	res, err := stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression)

	if err != nil {
		fmt.Println(task)
		fmt.Println("CreateTask(): breaky 3 ", err)
		return false, -1, err
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("CreateTask(): breaky 4 ", err)
		return false, -1, err
	}

	if task.IsRecurring {
		// rStmnt, err := tx.Preparex("INSERT INTO RecurrencePatterns (TaskID, RecurringType, DayOfWeek, DayOfMonth) VALUES (?, ?, ?, ?)")
		// if err != nil {
		// 	fmt.Println("CreateTask(): breaky 4", err)
		// 	return false, -1, err
		// }
		nexTime := cronexpr.MustParse("0 0 1 * * ?").NextN(task.StartTime, 10)
		fmt.Println(nexTime)
		rStmnt, err := tx.Preparex("INSERT INTO TaskTable (UserID, Category, TaskName, Description, StartTime, EndTime, Status, IsRecurring, IsAllDay, Difficulty, CronExpression) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println("CreateTask(): breaky 4", err)
			return false, -1, err
		}
		for i := 0; i < 5; i++ {
			_, err := rStmnt.Exec(task.UserID, task.Category, task.TaskName, task.Description, nexTime[i], nexTime[i].Add(task.EndTime.Sub(task.StartTime)), task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression)

			if err != nil {
				fmt.Println(task)
				fmt.Println("CreateTask(): breaky 7 ", err)
				return false, -1, err
			}
		}

	}
	// 	defer rStmnt.Close()

	// 	_, err = rStmnt.Exec(taskID, task.RecurringType, task.DayOfWeek, task.DayOfMonth)
	// 	if err != nil {
	// 		fmt.Println("CreateTask(): breaky 5", err)
	// 		return false, -1, err
	// 	}
	// }

	tx.Commit() //commit transaction to database
	fmt.Println("WE ADDED A TASK")
	return true, taskID, nil
}

func EditTask(task Task, id int) (bool, error) {

	tx, err := DB.Beginx()
	if err != nil {
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

	stmt1, err := tx.Preparex("DELETE FROM RecurringLog WHERE TaskID = ?")
	if err != nil {
		tx.Rollback()
		fmt.Println("Breaky; can't preparing statement for RecurringLog deletion:", err)
		return false, err
	}

	defer stmt1.Close()

	_, err = stmt1.Exec(id)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error deleting from RecurringLog:", err)
		return false, err
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

func Passtask(Tid int) bool {
	tx, err := DB.Beginx() // start transaction
	if err != nil {
		fmt.Printf("Passtask(): breaky 1 %v\n", err)
		return false
	}

	stmt, err := tx.Preparex(`
		UPDATE TaskTable 
		SET Status = ?
		WHERE TaskID = ?
	`)
	if err != nil {
		fmt.Printf("Passtask(): breaky 2 %v\n", err)
		tx.Rollback()
		return false
	}

	_, err = stmt.Exec("completed", Tid)
	if err != nil {
		fmt.Printf("Passtask(): breaky 3 %v\n", err)
		stmt.Close()
		tx.Rollback()
		return false
	}
	stmt.Close()
	tx.Commit()

	task, ok, err := GetTaskId(Tid)
	if err != nil {
		fmt.Printf("Passtask(): breaky 4 %v\n", err)
		return false
	}

	if !ok {
		fmt.Println("Passtask(): Task not found")
		return false
	}

	//tx, err = DB.Beginx() // start transaction

	points := CalculatePoints(task.Difficulty)
	_, err = DB.Exec("UPDATE UserTable SET Points = Points + ? WHERE UserID = ?", points, task.UserID)
	if err != nil {
		fmt.Printf("Passtask(): breaky 5 %v\n", err)
		tx.Rollback()
		return false
	}

	currBossHealth, _ := GetCurrBossHealth(task.UserID)

	// Check if the current boss health is zero
	if currBossHealth <= 0 {
		// Switch to the next boss ID (currBossId + 1)
		_, err := DB.Exec("UPDATE UserTable SET BossId = BossId + 1 WHERE UserID = ?", task.UserID)
		if err != nil {
			fmt.Printf("Passtask(): breaky 6 %v\n", err)
			return false
		}

		// Reset user points to 0
		_, err = DB.Exec("UPDATE UserTable SET Points = ? WHERE UserID = ?", 0, task.UserID)
		if err != nil {
			fmt.Printf("Passtask(): breaky 7 %v\n", err)
			return false
		}
	}

	//tx.Commit()
	return true
}

func Failtask(Tid int) bool {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		return false
	}

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

	return true

}
