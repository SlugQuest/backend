package crud

import (
	"log"
)

// Find user by UserID
func GetUser(uid string) (User, bool, error) {
	rows, err := DB.Query("SELECT * FROM UserTable WHERE UserID=?;", uid)
	var user User
	if err != nil {
		log.Println(err)
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

func GetUserPoints(Uid string) (int, bool, error) {
	log.Println(Uid)
	rows, err := DB.Query("SELECT Points FROM UserTable WHERE UserID = ?", Uid)
	thevalue := 0
	if err != nil {
		log.Println(err)
		log.Println("SOMETHING HAPPENED")
		rows.Close()
		return thevalue, false, err
	}
	counter := 0
	for rows.Next() {
		counter += 1
		log.Println(counter)
		rows.Scan(&thevalue)
		log.Println("finding")
	}
	rows.Close()

	return thevalue, counter == 1, err

}

// Add user into DB
func AddUser(u User) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("AddUser() #1: %v", err)
		return false, err
	}
	defer tx.Rollback() // aborrt transaction if error

	stmt, err := tx.Preparex("INSERT INTO UserTable (UserID, Points, Bossid) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("AddUser() #2: %v", err)
		return false, err
	}

	defer stmt.Close() //defer the closing of SQL statement to ensure it Closes once the function completes
	_, err = stmt.Exec(u.UserID, u.Points, u.BossId)
	if err != nil {
		log.Printf("AddUser() #3: %v", err)
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
	defer tx.Rollback() // aborrt transaction if error

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

func DeleteUser(uid string) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() // aborrt transaction if error

	stmt, err := tx.Preparex("DELETE FROM UserTable WHERE UserID = ?")
	if err != nil {
		log.Printf("DeleteUser() #1: %v", err)
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uid)
	if err != nil {
		log.Printf("DeleteUser() #2: %v", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}
