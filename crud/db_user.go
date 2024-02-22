package crud

import (
	"log"
	"math/rand"
	"time"
)

const USER_CODE_LEN int = 7

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
		rows.Scan(&user.UserID, &user.Points, &user.BossId, &user.SocialCode)
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
	socialCode, err := generateSocialCode()
	if err != nil {
		log.Printf("AddUser(): breaky 1: %v", err)
		return false, err
	}

	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("AddUser(): breaky 2: %v", err)
		return false, err
	}
	defer tx.Rollback() // abort transaction if error

	stmt, err := tx.Preparex("INSERT INTO UserTable (UserID, Points, Bossid, SocialCode) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("AddUser(): breaky 3: %v", err)
		return false, err
	}
	defer stmt.Close() //defer the closing of SQL statement to ensure it Closes once the function completes

	_, err = stmt.Exec(u.UserID, u.Points, u.BossId, socialCode)
	if err != nil {
		log.Printf("AddUser(): breaky 4: %v", err)
		return false, err
	}

	tx.Commit() //commit transaction to database

	return true, nil
}

// Generates a public code to differentiate users
func generateSocialCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codearr := make([]byte, USER_CODE_LEN)

	// Seed at the current time
	randgen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Loop until a unique ID was created
	isUnique := false
	var code string
	for !isUnique {
		// Generate a code
		for i := range codearr {
			codearr[i] = charset[randgen.Intn(len(charset))]
		}
		code = string(codearr)

		// No rows is desired in this case
		count := 0
		err := DB.Get(&count, "SELECT COUNT(*) FROM UserTable WHERE SocialCode = ?", code)
		if err != nil {
			return "", err
		}

		if count < 1 {
			isUnique = true
		}
	}

	return code, nil
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
		log.Println("DeleteUser: breaky 1")
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uid)
	if err != nil {
		log.Println("DeleteUser: breaky 2")
		return false, err
	}

	tx.Commit()

	return true, nil
}
