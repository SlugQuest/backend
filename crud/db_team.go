package crud

import (
	"fmt"
	"log"
	"time"

	"github.com/gorhill/cronexpr"
)


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

func AddUserToTeam(tid int64, ucode string) (bool, error) {
	if tid == int64(NoTeamID) {
		log.Println("AddUserToTeam(): invalid team")
		return false, nil
	}

	user, found, err := SearchUserCode(ucode, true)
	if !found || err != nil {
		log.Printf("AddUserToTeam() #1: could not find other user: %v", err)
		return false, err
	}
	uid, _ := user["UserID"].(string)
	prep, err := DB.Preparex("INSERT INTO TeamMembers (TeamID, UserID) VALUES (?,?)")
	if err != nil {
		log.Printf("AddUserToTeam(): could not prepare statement: %v", err)
		return false, err
	}
	_, err = prep.Exec(tid, uid)
	if err != nil {
		log.Printf("AddUserToTeam(): could not add team member: %v", err)
		return false, err
	}

	return true, nil
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
		log.Printf("GetUserTeams() #1: %v", err)
		rows.Close()
		return uteamArr, err
	}

	for rows.Next() {
		var taskprev Team
		err := rows.Scan(&taskprev.TeamID, &taskprev.Name)
		if err != nil {
			log.Printf("GetUserTeams(): could not read from DB: %v", err)
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
		log.Printf("GetTeamUsers() #1: %v", err)
		return users, err
	}
	rows, err := prep.Query(tid)
	if err != nil {
		log.Printf("GetTeamUsers() #2: %v", err)
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
			log.Printf("GetTeamUsers(): could not retreive users: %v", err)
			return users, err
		}
		log.Println("found a user", fUser)
		users = append(users, fUser)
	}
	return users, err

}

func RemoveUserFromTeam(tid int64, ucode string) (bool, error) {
	user, found, err := SearchUserCode(ucode, true)
	if !found || err != nil {
		log.Printf("RemoveUserFromTeam() #1: could not find other user: %v", err)
		return false, err
	}
	uid, _ := user["UserID"].(string)
	prep, err := DB.Preparex("DELETE FROM TeamMembers WHERE TeamID = ? AND UserID = ?")
	if err != nil {
		log.Printf("RemoveUserFromTeam() #2: could not prepare statement: %v", err)
		return false, err
	}
	_, err = prep.Exec(tid, uid)
	if err != nil {
		log.Printf("RemoveUserFromTeam() #3: could not remove team member: %v", err)
		return false, err
	}
	return true, err
}

func DeleteTeam(tid int64) (bool, error) {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		log.Printf("DeleteTeam(): DB issue starting transaction: %v", err)
		return false, err
	}
	defer tx.Rollback()

	stmnt, err := tx.Preparex("DELETE FROM TeamMembers WHERE TeamID = ? ")
	if err != nil {
		log.Printf("DeleteTeam() #1: could not prepare statement: %v", err)
		return false, err
	}
	_, err = stmnt.Exec(tid)
	if err != nil {
		log.Printf("DeleteTeam() #2: could not delete team members: %v", err)
		return false, err
	}

	stmnt2, err := tx.Preparex("DELETE FROM TeamTable WHERE TeamID = ?")
	if err != nil {
		log.Printf("DeleteTeam() #3: could not prepare statement: %v", err)
		return false, err
	}

	_, err = stmnt2.Exec(tid)
	if err != nil {
		log.Printf("DeleteTeam() #4: could not delete team: %v", err)
		return false, err
	}
	tx.Commit()
	return true, nil
}

func CreateTeam(name string, uid string) (bool, int64, error) {
	tx, err := DB.Beginx() //start transaction
	if err != nil {
		log.Printf("CreateTeam(): DB issue starting transaction: %v", err)
		return false, 0, err
	}
	defer tx.Rollback()

	stmnt, err := tx.Preparex("INSERT INTO TeamTable (TeamName) VALUES (?)")
	if err != nil {
		log.Printf("CreateTeam(): could not prepare statement: %v", err)
		return false, 0, err
	}
	res, err := stmnt.Exec(name)
	if err != nil {
		log.Printf("CreateTeam(): could not create team: %v", err)
		return false, 0, err
	}
	teamins, err := res.LastInsertId()
	if err != nil {
		log.Printf("CreateTeam(): breaky 3: %v", err)
		return false, 0, err
	}
	tx.Commit()
	AddUserToTeamUid(teamins, uid)

	return true, teamins, nil
}

func AddUserToTeamUid(tid int64, uid string) (bool, error) {
	prep, err := DB.Preparex("INSERT INTO TeamMembers (TeamID, UserID) VALUES (?,?)")
	if err != nil {
		log.Printf("AddUserToTeamUid(): could not prepare statement %v", err)
		return false, err
	}
	_, err = prep.Exec(tid, uid)
	if err != nil {
		log.Printf("AddUserToTeamUid(): could not add user to team: %v", err)
		return false, err
	}

	return true, nil
}