package crud

import (
	"fmt"
	"log"
	"time"

	"github.com/gorhill/cronexpr"
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
func GetUserTask(uid string) ([]Task, error) {
	utaskArr := []Task{}

	prep, err := DB.Preparex(`SELECT TaskID, UserID, Category, TaskName, Description, StartTime, EndTime, Status, IsRecurring, IsAllDay, Difficulty, CronExpression FROM TaskTable
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
		var taskprev Task
		err := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.Description, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay, &taskprev.Difficulty, &taskprev.CronExpression)
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

func GetTeamTask(tid int) ([]Task, error) {
	utaskArr := []Task{}

	prep, err := DB.Preparex(`SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, Status, IsRecurring, IsAllDay, Description, Difficulty FROM TaskTable
		WHERE TeamID = ?;`)
	if err != nil {
		log.Printf("GetTeamTask() #1: %v", err)
		return utaskArr, err
	}

	rows, err := prep.Query(tid)
	if err != nil {
		log.Printf("GetTeamTask() #2: %v", err)
		rows.Close()
		prep.Close()
		return utaskArr, err
	}

	for rows.Next() {
		var taskprev Task
		err := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay, &taskprev.Description, &taskprev.Difficulty)
		if err != nil {
			log.Printf("GetTeamTask() #3: %v", err)
			rows.Close()
		}
		utaskArr = append(utaskArr, taskprev)
	}
	prep.Close()
	rows.Close()
	return utaskArr, err
}

func AddUserToTeam(tid int64, ucode string) bool {
	user, found, err := SearchUserCode(ucode, true)
	if !found || err != nil {
		log.Printf("AddFriend() #1: could not find other user: %v", err)
		return false
	}
	uid, _ := user["UserID"].(string)
	prep, err := DB.Preparex("INSERT INTO TeamMembers (TeamID, UserID) VALUES (?,?)")
	if err != nil {
		log.Printf("bricked in adduser team")
		return false
	}
	_, err = prep.Exec(tid, uid)
	if err != nil {
		log.Printf("bricked in adduser team")
		return false
	}
	return true

}

func GetUserTeams(uid string) ([]Team, error) {
	uteamArr := []Team{}

	prep, err := DB.Preparex("SELECT t.TeamID, t.TeamName FROM TeamMembers z, TeamTable t WHERE UserID = ? AND t.TeamID = z.TeamID ")
	if err != nil {
		log.Println(err)
		return uteamArr, err
	}
	rows, err := prep.Query(uid)
	if err != nil {
		log.Println(err)
		return uteamArr, err
	}
	if err != nil {
		log.Printf("getuserteamissue")
		rows.Close()
		return uteamArr, err
	}

	for rows.Next() {
		var taskprev Team
		err := rows.Scan(&taskprev.TeamID, &taskprev.Name)
		if err != nil {
			fmt.Println(err)
			rows.Close()
			return uteamArr, err
		}
		log.Println("we found a team")
		taskprev.Members, _ = GetTeamUsers(taskprev.TeamID)
		uteamArr = append(uteamArr, taskprev)
	}
	
	return uteamArr, nil

}

func GetTeamUsers(tid int64) ([]map[string]interface{}, error) {
	uarr := []string{}
	var users []map[string]interface{}
	prep, err := DB.Preparex("SELECT u.UserID FROM UserTable u, TeamMembers m WHERE u.UserID = m.UserID AND m.TeamID = ?")
	if err != nil {
		log.Printf("getuserteamissue", err)
		return users, err
	}
	rows, err := prep.Query(tid)
	if err != nil {
		log.Printf("getuserteamissue", err)
		rows.Close()
		return users, err
	}

	for rows.Next() {
		var uid string
		log.Println("found a user")
		err := rows.Scan(&uid)
		if err != nil {
			fmt.Println(err)
			rows.Close()
			return users, err
		}
		uarr = append(uarr, uid)
	}
	rows.Close()
	for _, fID := range uarr {
		fUser, found, err := GetPublicUser(fID)
		if !found || err != nil {
			log.Printf("Get TEam useres could not retreive users  %v", err)
			return users, err
		}
		log.Println("found a user", fUser)
		users = append(users, fUser)
	}
	return users, err

}

func RemoveUserFromTeam(tid int64, ucode string) bool {
	user, found, err := SearchUserCode(ucode, true)
	if !found || err != nil {
		log.Printf("AddFriend() #1: could not find other user: %v", err)
		return false
	}
	uid, _ := user["UserID"].(string)
	prep, err := DB.Preparex("DELETE FROM TeamMembers WHERE TeamID = ? AND UserID = ?")
	if err != nil {
		log.Printf("bricked in del team")
		return false
	}
	_, err = prep.Exec(tid, uid)
	if err != nil {
		log.Printf("bricked in delteam")
		return false
	}
	return true

}

func DeleteTeam(tid int64) bool {
	tx, err := DB.Beginx() //start transaction
	defer tx.Rollback()
	stmnt, err := tx.Preparex("DELETE FROM TeamMembers WHERE TeamID = ? ")
	if err != nil {
		log.Printf("bricked in del team")
		return false
	}
	_, err = stmnt.Exec(tid)
	if err != nil {
		log.Printf("bricked in del team")
		return false
	}

	stmnt2, err := tx.Preparex("DELETE FROM TeamTable WHERE TeamID = ?")
	if err != nil {
		log.Printf("bricked in del team")
		return false
	}

	_, err = stmnt2.Exec(tid)
	if err != nil {
		log.Printf("bricked in del team")
		return false
	}
	tx.Commit()
	return true
}

func CreateTeam(name string, uid string) (bool, int64) {
	tx, err := DB.Beginx() //start transaction
	defer tx.Rollback()
	stmnt, err := tx.Preparex("INSERT INTO TeamTable (TeamName) VALUES (?)")
	if err != nil {
		log.Printf("bricked in add team", err)
		return false, 0
	}
	res, err := stmnt.Exec(name)
	if err != nil {
		log.Printf("bricked in add team", err)
		return false, 0
	}
	teamins, err := res.LastInsertId()
	if err != nil {
		// fmt.Println(task)
		fmt.Println("CreateTeam(): breaky 3 ", err)
		return false, 0
	}
	tx.Commit()
	AddUserToTeamUid(teamins, uid)

	return true, teamins
}

func AddUserToTeamUid(tid int64, uid string) bool {
	prep, err := DB.Preparex("INSERT INTO TeamMembers (TeamID, UserID) VALUES (?,?)")
	if err != nil {
		log.Printf("bricked in adduser team", err)
		return false
	}
	_, err = prep.Exec(tid, uid)
	if err != nil {
		log.Printf("bricked in adduser team", err)
		return false
	}
	return true

}

func GetUserTaskDateTime(uid string, startq time.Time, endq time.Time) ([]RecurTypeTask, error) {
	utaskArr := []RecurTypeTask{}

	prep, err := DB.Preparex(`SELECT TaskID, UserID, Category, TaskName, StartTime, EndTime, Status, IsRecurring, IsAllDay FROM TaskTable t 
		WHERE UserID = ? AND t.StartTime > ? AND t.EndTime < ? AND NOT IsRecurring;`)
	if err != nil {
		log.Printf("GetUserTaskDateTime() #1: %v", err)
		return utaskArr, err
	}
	defer prep.Close()
	log.Printf("GetUserTaskDateTime() #2: %v", err)
	log.Println(startq)
	rows, err := prep.Query(uid, startq, endq)
	if err != nil {
		log.Printf("GetUserTaskDateTime() #2: %v", err)
		rows.Close()
		return utaskArr, err
	}

	for rows.Next() {
		var taskprev RecurTypeTask
		err := rows.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay)
		log.Println("found not recur")
		if err != nil {
			log.Printf("GetUserTaskDateTime() #3: %v", err)
			rows.Close()
			return utaskArr, err
		}
		taskprev.RecurrenceId = 1
		utaskArr = append(utaskArr, taskprev)
	}
	prep.Close()
	rows.Close()
	p2, err := DB.Preparex("SELECT c.TaskID, c.UserID, c.Category, c.TaskName, c.StartTime, c.EndTime, c.Status, c.IsRecurring, c.IsAllDay, l.timestamp, l.LogId FROM TaskTable c, RecurringLog l WHERE l.TaskID = c.TaskID AND  c.UserID = ? AND l.timestamp > ? AND l.timestamp < ?;")
	if err != nil {
		log.Println("found recur")
		log.Printf("GetUserTaskDateTime() #4: %v", err)
		return utaskArr, err
	}
	rowrec, err := p2.Query(uid, startq, endq)
	for rows.Next() {
		var taskprev RecurTypeTask
		log.Println("recur type task in span found")
		var reftime time.Time
		err := rowrec.Scan(&taskprev.TaskID, &taskprev.UserID, &taskprev.Category, &taskprev.TaskName, &taskprev.StartTime, &taskprev.EndTime, &taskprev.Status, &taskprev.IsRecurring, &taskprev.IsAllDay, &reftime, &taskprev.RecurrenceId)
		if err != nil {
			log.Printf("GetUserTaskDateTime() #5: %v", err)
			rowrec.Close()
			return utaskArr, err
		}
		taskprev.StartTime = taskprev.StartTime
		utaskArr = append(utaskArr, taskprev)
	}
	p2.Close()
	rowrec.Close()
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

	tx.Commit() //commit transaction to database

	if task.IsRecurring {
		currentMonth := time.Now().Month()
		currentYear := time.Now().Year()
		nextTimes := cronexpr.MustParse(task.CronExpression).NextN(time.Now(), 31)
		//assuming there can only be one recurrence a day, so at most 31 recurrences in a month

		for _, nextTime := range nextTimes {
			// Check if the next occurrence is in the current month
			if nextTime.Month() == currentMonth && nextTime.Year() == currentYear {
				_, _, err = CreateRecurringLogEntry(task.TaskID, "todo", nextTime)
				if err != nil {
					fmt.Printf("In here")
					return false, -1, err
				}
			}
		}
	}

	return true, taskID, nil
}

func EditTask(task Task, tid int) (bool, error) {
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
	log.Printf("thetaskis bieng editedis")
	log.Printf(task.TaskName)
	_, err = stmt.Exec(task.Category, task.TaskName, task.Description, task.StartTime, task.EndTime, task.Status, task.IsRecurring, task.IsAllDay, task.Difficulty, task.CronExpression,
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

func Passtask(Tid int, uid string) (bool, error, int) {
	task, ok, err := GetTaskId(Tid)
	if err != nil {
		fmt.Printf("Passtask(): breaky 2 %v\n", err)
		return false, err, -1
	}

	if !ok {
		fmt.Println("Passtask(): Task not found")
		return false, fmt.Errorf("task not found"), -1
	}

	if task.UserID != uid {
		return false, fmt.Errorf("task not owned by this user"), -1
	}

	if task.IsRecurring {
		_, err := DB.Exec(`
		UPDATE RecurringLog 
		SET Status = ?
		WHERE (LogId, timestamp) in (SELECT LogId, timestamp from (SELECT r.LogId, MIN(r.timestamp) FROM TaskTable t, RecurringLog r WHERE t.TaskID = r.TaskID AND t.TaskID = ? AND r.timestamp > ?) as temptable)
	`, "completed", Tid, time.Now())

		if err != nil {
			fmt.Printf("Passtask(): breaky 0 %v\n", err)
			return false, err, -1
		}

	} else {
		tx, err := DB.Beginx() // start transaction
		if err != nil {
			fmt.Printf("Passtask(): breaky 1 %v\n", err)
			return false, err, -1
		}
		defer tx.Rollback() // Abort transaction if any error occurs
		stmt, err := tx.Preparex(`
			UPDATE TaskTable 
			SET Status = ?
			WHERE TaskID = ?
		`)

		if err != nil {
			fmt.Printf("Passtask(): breaky 2 %v\n", err)
			return false, err, -1
		}

		_, err = stmt.Exec("completed", Tid)
		if err != nil {
			fmt.Printf("Passtask(): breaky 3 %v\n", err)
			return false, err, -1
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
		return false, err, -1
	}

	currBossHealth, err := GetCurrBossHealth(task.UserID)
	if err != nil {
		fmt.Printf("Passtask(): breaky %v\n", err)
		return false, err, -1
	}

	// Check if the current boss health is zero
	if currBossHealth <= 0 {
		// Switch to the next boss ID (currBossId + 1)
		_, err := DB.Exec("UPDATE UserTable SET BossId = BossId + 1 WHERE UserID = ?", task.UserID)
		if err != nil {
			fmt.Printf("Passtask(): breaky 6 %v\n", err)
			return false, err, -1
		}

		// Reset user points to 0
		_, err = DB.Exec("UPDATE UserTable SET Points = ? WHERE UserID = ?", 0, task.UserID)
		if err != nil {
			fmt.Printf("Passtask(): breaky 7 %v\n", err)
			return false, err, -1
		}
	}

	user, _, _ := GetUser(uid)

	return true, nil, user.BossId
}

func PassRecurringTask(Tid int, recurrenceID int, uid string) (bool, error, int) {
	task, ok, err := GetTaskId(Tid)
	if err != nil {
		fmt.Printf("Passtask(): breaky 0 %v\n", err)
		return false, err, -1
	}

	if !ok {
		fmt.Println("Passtask(): Task not found")
		return false, fmt.Errorf("task not found"), -1
	}

	if task.UserID != uid {
		return false, fmt.Errorf("task not owned by this user"), -1
	}

	tx, err := DB.Beginx()
	if err != nil {
		fmt.Printf("PassRecurringTask(): breaky 1 %v\n", err)
		return false, err, -1
	}
	defer tx.Rollback()

	_, err = DB.Exec(`
		UPDATE RecurringLog 
		SET Status = ?
		WHERE LogId = ?
	`, "completed", recurrenceID)

	if err != nil {
		fmt.Printf("PassRecurringTask(): breaky 2 %v\n", err)
		return false, err, -1
	}

	tx.Commit()

	// Update user points
	points := CalculatePoints(task.Difficulty)
	_, err = DB.Exec("UPDATE UserTable SET Points = Points + ? WHERE UserID = ?", points, task.UserID)
	if err != nil {
		fmt.Printf("PassRecurringTask(): breaky 3 %v\n", err)
		return false, err, -1
	}

	currBossHealth, err := GetCurrBossHealth(uid)
	if err != nil {
		fmt.Printf("PassRecurringTask(): breaky 4 %v\n", err)
		return false, err, -1
	}

	// Check if the current boss health is zero
	if currBossHealth <= 0 {
		// Switch to the next boss ID (currBossId + 1)
		_, err := DB.Exec("UPDATE UserTable SET BossId = BossId + 1 WHERE UserID = ?", task.UserID)
		if err != nil {
			fmt.Printf("PassRecurringTask(): breaky 5 %v\n", err)
			return false, err, -1
		}

		// Reset user points to 0
		_, err = DB.Exec("UPDATE UserTable SET Points = ? WHERE UserID = ?", 0, task.UserID)
		if err != nil {
			fmt.Printf("PassRecurringTask(): breaky 6 %v\n", err)
			return false, err, -1
		}
	}

	user, _, _ := GetUser(uid)

	return true, nil, user.BossId
}

func Failtask(Tid int, uid string) (bool, error) {
	task, ok, err := GetTaskId(Tid)
	if err != nil {
		fmt.Printf("Failtask(): breaky %v\n", err)
		return false, err
	}

	if !ok {
		fmt.Println("Failtask(): Task not found")
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
	`, "failed", Tid, time.Now())

		if err != nil {
			fmt.Printf("Failtask(): breaky 0 %v\n", err)
			return false, err
		}
	} else {
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
			return false, err
		}
		swag, err := stmt.Exec("failed", Tid)
		stmt.Close()
		if err != nil {
			print(err.Error())
			print("FailtTask(): breaky 1 ")
			fmt.Println(err)
			fmt.Println(swag)
			return false, err
		}

		tx.Commit()
	}

	return true, nil
}

func FailRecurringTask(Tid int, recurrenceID int, uid string) (bool, error) {
	task, ok, err := GetTaskId(Tid)
	if err != nil {
		fmt.Printf("FailRecurringTask(): breaky %v\n", err)
		return false, err
	}

	if !ok {
		fmt.Println("FailRecurringTask(): Task not found")
		return false, fmt.Errorf("task not found")
	}

	if task.UserID != uid {
		return false, fmt.Errorf("task not owned by this user")
	}

	tx, err := DB.Beginx()
	if err != nil {
		fmt.Printf("FailRecurringTask(): breaky 1 %v\n", err)
		return false, err
	}
	defer tx.Rollback()

	_, err = DB.Exec(`
		UPDATE RecurringLog 
		SET Status = ?
		WHERE LogId = ?
	`, "failed", recurrenceID)

	if err != nil {
		fmt.Printf("FailRecurringTask(): breaky 2 %v\n", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}
