package crud

import (
	"fmt"
	"log"
	"time"

	"github.com/gorhill/cronexpr"
)

// GetTaskId finds a task by TaskID.
// Input: tid (int) - TaskID
// Output: Task - Retrieved task, bool - Success flag, error - Potential error
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
		err := rows.Scan(&taskit.TaskID, &taskit.UserID, &taskit.Category, &taskit.TaskName, &taskit.Description, &taskit.StartTime, &taskit.EndTime, &taskit.Status, &taskit.IsRecurring, &taskit.IsAllDay, &taskit.Difficulty, &taskit.CronExpression, &taskit.TeamID)
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
// GetUserTask retrieves tasks associated with a specific user.
// Input: uid (string) - UserID
// Output: []Task - List of tasks, error - Potential error
func GetUserTask(uid string) ([]Task, error) {
	utaskArr := []Task{}

	prep, err := DB.Preparex(`SELECT TaskID, UserID, Category, TaskName, Description, StartTime, EndTime, Status, IsRecurring, IsAllDay, Difficulty, CronExpression, TeamID
		FROM TaskTable WHERE UserID = ?;`)
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
		var taskprev Task
		err := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.Description, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay, &taskprev.Difficulty, &taskprev.CronExpression, &taskprev.TeamID)
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

// GetUserTaskDateTime retrieves tasks for a user within a specific time range.
// Input: uid (string) - UserID, startq (time.Time) - Start time, endq (time.Time) - End time
// Output: []RecurTypeTask - List of tasks, error - Potential error
func GetUserTaskDateTime(uid string, startq time.Time, endq time.Time) ([]RecurTypeTask, error) {
	utaskArr := []RecurTypeTask{}

	prep, err := DB.Preparex(`SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, Status, IsRecurring, IsAllDay FROM TaskTable t 
		WHERE UserID = ? AND t.StartTime > ? AND t.EndTime < ? AND NOT IsRecurring;`)
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
		var taskprev RecurTypeTask
		err := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay)
		if err != nil {
			log.Printf("GetUserTaskDateTime() #3: %v", err)
			rows.Close()
			return utaskArr, err
		}

		taskprev.RecurrenceId = -1
		utaskArr = append(utaskArr, taskprev)
	}
	prep.Close()
	rows.Close()

	p2, err := DB.Preparex("SELECT c.TaskID, c.UserID, c.Category, c.TaskName, c.StartTime, c.EndTime, c.Status, c.IsRecurring, c.IsAllDay, l.timestamp, l.LogId FROM TaskTable c, RecurringLog l WHERE l.TaskID = c.TaskID AND c.UserID = ? AND l.timestamp > ? AND l.timestamp < ?;")
	if err != nil {
		log.Printf("GetUserTaskDateTime() #4: %v", err)
		return utaskArr, err
	}
	rowrec, err := p2.Query(uid, startq, endq)
	for rowrec.Next() {
		var taskprev RecurTypeTask
		log.Println("recur type task in span found")
		var reftime time.Time
		err := rowrec.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay, &reftime, &taskprev.RecurrenceId)
		if err != nil {
			log.Printf("GetUserTaskDateTime() #5: %v", err)
			rowrec.Close()
			return utaskArr, err
		}

		taskprev.EndTime = reftime.Add(taskprev.EndTime.Sub(taskprev.StartTime))
		taskprev.StartTime = reftime
		utaskArr = append(utaskArr, taskprev)
	}
	p2.Close()
	rowrec.Close()
	return utaskArr, err
}

// CreateTask creates a new task and handles recurring tasks.
// Input: task (Task) - Task details
// Output: bool - Success flag, int64 - TaskID, error - Potential error
func CreateTask(task Task) (bool, int64, error) {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		log.Printf("CreateTask(): DB issue starting transaction: %v", err)
		return false, -1, err
	}
	defer tx.Rollback() // Abort transaction if any error occurs

	//preparing statement to prevent SQL injection issues
	stmt, err := tx.Preparex("INSERT INTO TaskTable (UserID, Category, TaskName, Description, StartTime, EndTime, Status, IsRecurring, IsAllDay, Difficulty, CronExpression, TeamID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("CreateTask(): could not prepare statement %v", err)
		return false, -1, err
	}
	defer stmt.Close() // Defer the closing of SQL statement to ensure it closes once the function completes

	res, err := stmt.Exec(task.UserID, task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression, task.TeamID)
	if err != nil {
		log.Printf("CreateTask(): could not insert into table: %v", err)
		return false, -1, err
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		log.Printf("CreateTask(): breaky 4: %v", err)
		return false, -1, err
	}

	tx.Commit() //commit transaction to database

	if task.IsRecurring {
		currentTime := time.Now()
		currentMonth := currentTime.Month()
		currentYear := currentTime.Year()
		nextTimes := cronexpr.MustParse(task.CronExpression).NextN(time.Now(), 31)
		//assuming there can only be one recurrence a day, so at most 31 recurrences in a month

		if task.StartTime.Day() == currentTime.Day() &&
			task.StartTime.Month() == currentMonth &&
			task.StartTime.Year() == currentYear {
			// currentTime.After(task.StartTime) {
			_, _, err := CreateRecurringLogEntry(taskID, "todo", task.StartTime)
			if err != nil {
				return false, -1, err
			}
		}

		for _, nextTime := range nextTimes {
			// Check if the next occurrence is in the current month
			if nextTime.Month() == currentMonth && nextTime.Year() == currentYear {
				_, _, err = CreateRecurringLogEntry(taskID, "todo", nextTime)
				if err != nil {
					return false, -1, err
				}
			}
		}
	}

	return true, taskID, nil
}

// EditTask updates an existing task.
// Input: task (Task) - Task details, tid (int) - TaskID
// Output: bool - Success flag, error - Potential error
func EditTask(task Task, tid int) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("EditTask() #1: DB issue starting transaction: %v", err)
		return false, err
	}

	stmt, err := tx.Preparex(`
		UPDATE TaskTable 
		SET Category = ?, TaskName = ?, Description = ?, StartTime = ?, EndTime = ?, Status = ?, IsRecurring = ?, IsAllDay = ?, Difficulty = ?, CronExpression = ?, TeamID = ? 
		WHERE TaskID = ? AND UserID = ?
	`)

	if err != nil {
		log.Printf("EditTask() #2: could not prepare statement: %v", err)
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression, task.TeamID,
		tid, task.UserID)
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

// DeleteTask deletes a task and its associated recurring logs.
// Input: tid (int) - TaskID, uid (string) - UserID
// Output: bool - Success flag, error - Potential error
func DeleteTask(tid int, uid string) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("DeleteTask() #1: %v", err)
		return false, err
	}
	defer tx.Rollback() // Abort transaction if any error occurs

	delTT, err := tx.Preparex("DELETE FROM TaskTable WHERE TaskID = ? AND UserID = ?")
	if err != nil {
		log.Printf("DeleteTask() #2: %v", err)
		return false, err
	}
	defer delTT.Close()

	res, err := delTT.Exec(tid, uid)
	if err != nil {
		log.Printf("DeleteTask() #3: %v", err)
		return false, err
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		log.Printf("DeleteTask() #4: %v", err)
		return false, err
	}

	// Only delete from RecurringLog (not valited by UID) if any were deleted from the main table
	if numDeleted > 0 {
		delRL, err := tx.Preparex("DELETE FROM RecurringLog WHERE TaskID = ?")
		if err != nil {
			log.Printf("DeleteTask() #5: can't preparing statement for RecurringLog deletion: %v", err)
			return false, err
		}
		defer delRL.Close()

		_, err = delRL.Exec(tid)
		if err != nil {
			log.Printf("DeleteTask() #6: Error deleting from RecurringLog: %v", err)
			return false, err
		}
	}

	tx.Commit()

	return true, nil
}

// PassTask marks a task as completed, updates user points, and potentially switches to the next boss.
// Input: tid (int) - TaskID, uid (string) - UserID
// Output: bool - Success flag, int - Next boss ID, error - Potential error
func PassTask(tid int, uid string) (bool, int, error) {
	task, found, err := GetTaskId(tid)
	if err != nil {
		log.Printf("PassTask(): error getting task: %v\n", err)
		return false, -1, err
	}
	if !found {
		log.Println("PassTask(): Task not found")
		return false, -1, fmt.Errorf("task not found")
	}

	if task.UserID != uid {
		return false, -1, fmt.Errorf("task not owned by this user")
	}

	if task.IsRecurring {
		_, err := DB.Exec(`
			UPDATE RecurringLog 
			SET Status = ?
			WHERE (LogId, timestamp) in (SELECT LogId, timestamp from (SELECT r.LogId, MIN(r.timestamp) FROM TaskTable t, RecurringLog r WHERE t.TaskID = r.TaskID AND t.TaskID = ? AND r.timestamp > ?) as temptable)
		`, "completed", tid, time.Now())

		if err != nil {
			log.Printf("PassTask(): error updating RecurringLog: %v\n", err)
			return false, -1, err
		}

	} else {
		tx, err := DB.Beginx() // start transaction
		if err != nil {
			log.Printf("PassTask(): DB error starting transaction: %v\n", err)
			return false, -1, err
		}
		defer tx.Rollback() // Abort transaction if any error occurs

		stmt, err := tx.Preparex(`
			UPDATE TaskTable 
			SET Status = ?
			WHERE TaskID = ?
		`)
		if err != nil {
			log.Printf("PassTask(): could not prepare statement: %v\n", err)
			return false, -1, err
		}
		defer stmt.Close()

		_, err = stmt.Exec("completed", tid)
		if err != nil {
			log.Printf("PassTask(): could not update status: %v\n", err)
			return false, -1, err
		}

		tx.Commit()
	}

	points := CalculatePoints(task.Difficulty)
	_, err = DB.Exec("UPDATE UserTable SET Points = Points + ? WHERE UserID = ?", points, task.UserID)
	if err != nil {
		fmt.Printf("PassTask(): could not update user's points: %v\n", err)
		return false, -1, err
	}

	currBossHealth, err := GetCurrBossHealth(task.UserID)
	if err != nil {
		fmt.Printf("PassTask(): could not retrieve current boss health: %v\n", err)
		return false, -1, err
	}

	// Check if the current boss health is zero
	if currBossHealth <= 0 {
		// Switch to the next boss ID (currBossId + 1)
		_, err := DB.Exec("UPDATE UserTable SET BossId = BossId + 1 WHERE UserID = ?", task.UserID)
		if err != nil {
			fmt.Printf("PassTask(): could not switch to next boss: %v\n", err)
			return false, -1, err
		}

		// Reset user points to 0
		_, err = DB.Exec("UPDATE UserTable SET Points = ? WHERE UserID = ?", 0, task.UserID)
		if err != nil {
			fmt.Printf("PassTask(): could not reset user points to 0: %v\n", err)
			return false, -1, err
		}
	}

	user, _, _ := GetUser(uid)
	return true, user.BossId, nil
}

// PassRecurringTask marks a recurring task as completed, updates user points, and potentially switches to the next boss.
// Input: tid (int) - TaskID, recurrenceID (int) - Recurring Log ID, uid (string) - UserID
// Output: bool - Success flag, int - user boss ID, error - Potential error
func PassRecurringTask(tid int, recurrenceID int, uid string) (bool, int, error) {
	task, ok, err := GetTaskId(tid)
	if err != nil {
		log.Printf("PassRecurringTask(): error getting task: %v\n", err)
		return false, -1, err
	}

	if !ok {
		fmt.Println("PassTask(): Task not found")
		return false, -1, fmt.Errorf("task not found")
	}

	if task.UserID != uid {
		return false, -1, fmt.Errorf("task not owned by this user")
	}

	tx, err := DB.Beginx()
	if err != nil {
		fmt.Printf("PassRecurringTask(): DB issue starting transaction %v\n", err)
		return false, -1, err
	}
	defer tx.Rollback()

	_, err = DB.Exec(`
		UPDATE RecurringLog 
		SET Status = ?
		WHERE LogId = ?
	`, "completed", recurrenceID)

	if err != nil {
		fmt.Printf("PassRecurringTask(): could not update status: %v\n", err)
		return false, -1, err
	}

	tx.Commit()

	// Update user points
	points := CalculatePoints(task.Difficulty)
	_, err = DB.Exec("UPDATE UserTable SET Points = Points + ? WHERE UserID = ?", points, task.UserID)
	if err != nil {
		fmt.Printf("PassRecurringTask(): could not update user points: %v\n", err)
		return false, -1, err
	}

	currBossHealth, err := GetCurrBossHealth(uid)
	if err != nil {
		fmt.Printf("PassRecurringTask(): could not retrieve current boss health: %v\n", err)
		return false, -1, err
	}

	// Check if the current boss health is zero
	if currBossHealth <= 0 {
		// Switch to the next boss ID (currBossId + 1)
		_, err := DB.Exec("UPDATE UserTable SET BossId = BossId + 1 WHERE UserID = ?", task.UserID)
		if err != nil {
			fmt.Printf("PassRecurringTask(): could not switch to next boss %v\n", err)
			return false, -1, err
		}

		// Reset user points to 0
		_, err = DB.Exec("UPDATE UserTable SET Points = ? WHERE UserID = ?", 0, task.UserID)
		if err != nil {
			fmt.Printf("PassRecurringTask(): could not reset user points to 0 %v\n", err)
			return false, -1, err
		}
	}

	user, _, _ := GetUser(uid)

	return true, user.BossId, nil
}

// FailTask marks a task as failed, subtracts points if it was previously completed.
// Input: tid (int) - TaskID, uid (string) - UserID
// Output: bool - Success flag, error - Potential error
func FailTask(tid int, uid string) (bool, error) {
	task, found, err := GetTaskId(tid)
	if err != nil {
		log.Printf("FailTask(): error getting task: %v\n", err)
		return false, err
	}

	if !found {
		log.Println("FailTask(): Task not found")
		return false, fmt.Errorf("task not found")
	}

	if task.UserID != uid {
		return false, fmt.Errorf("task not owned by this user")
	}

	if task.IsRecurring {
		_, err := DB.Exec(`
			UPDATE RecurringLog 
			SET Status = ?
			WHERE (LogId, timestamp) in (SELECT LogId, timestamp from (SELECT r.LogId, MIN(r.timestamp) FROM TaskTable t, RecurringLog r WHERE t.TaskID = r.TaskID AND t.TaskID = ? AND r.timestamp > ?) as temptable)
		`, "failed", tid, time.Now())

		if err != nil {
			fmt.Printf("FailTask(): error updating RecurringLog: %v\n", err)
			return false, err
		}
	} else {
		passedBefore := task.Status == "completed"
		tx, err := DB.Beginx() //start transaction
		if err != nil {
			return false, err
		}
		defer tx.Rollback() // Abort transaction if any error occurs

		stmt, err := tx.Preparex(`
			UPDATE TaskTable 
			SET Status = ?
			WHERE TaskID = ?
		`)
		if err != nil {
			log.Printf("FailTask(): could not prepare statement: %v\n", err)
			return false, err
		}
		defer stmt.Close()

		_, err = stmt.Exec("failed", tid)
		if err != nil {
			log.Printf("FailTask(): could not update status: %v\n", err)
			return false, err
		}

		tx.Commit()

		if passedBefore {
			points := CalculatePoints(task.Difficulty)
			_, err = DB.Exec("UPDATE UserTable SET Points = Points - ? WHERE UserID = ?", points, task.UserID)
			if err != nil {
				log.Printf("PassTask(): could not update user's points: %v\n", err)
				return false, err
			}
		}
	}

	return true, nil
}

// FailRecurringTask marks a recurring task as failed.
// Input: tid (int) - TaskID, recurrenceID (int) - Recurring Log ID, uid (string) - UserID
// Output: bool - Success flag, error - Potential error
func FailRecurringTask(tid int, recurrenceID int, uid string) (bool, error) {
	task, ok, err := GetTaskId(tid)
	if err != nil {
		log.Printf("FailRecurringTask(): error getting task %v\n", err)
		return false, err
	}

	if !ok {
		log.Println("FailRecurringTask(): Task not found")
		return false, fmt.Errorf("task not found")
	}

	if task.UserID != uid {
		return false, fmt.Errorf("task not owned by this user")
	}

	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("FailRecurringTask(): DB issue starting transaction: %v", err)
		return false, err
	}
	defer tx.Rollback()

	_, err = DB.Exec(`
		UPDATE RecurringLog 
		SET Status = ?
		WHERE LogId = ?
	`, "failed", recurrenceID)

	if err != nil {
		log.Printf("FailRecurringTask(): could not update status: %v\n", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}
