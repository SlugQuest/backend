package crud

import (
	"log"
	"math/rand"
	"time"
)

const socialcode_set = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const USER_CODE_LEN int = 7

// Find user by UserID
func GetUser(uid string) (User, bool, error) {
	rows, err := DB.Query("SELECT * FROM UserTable WHERE UserID=?;", uid)
	var user User
	if err != nil {
		log.Printf("GetUser() #1: %v", err)
		return user, false, err
	}

	counter := 0
	for rows.Next() {
		counter += 1
		rows.Scan(&user.UserID, &user.Username, &user.Picture, &user.Points, &user.BossId, &user.SocialCode)
	}
	rows.Close()

	return user, counter == 1, err
}

// Return public information about a user
func GetPublicUser(uid string) (map[string]interface{}, bool, error) {
	user, found, err := GetUser(uid)
	if err != nil {
		log.Printf("GetPublicUser() #1: %v", err)
		return map[string]interface{}{}, false, err
	}

	if !found {
		log.Println("GetPublicUser() #2: did not find user")
		return map[string]interface{}{}, false, err
	}

	publicUser := map[string]interface{}{
		"Username":   user.Username,
		"Picture":    user.Picture,
		"Points":     user.Points,
		"BossId":     user.BossId,
		"SoicalCode": user.SocialCode,
	}

	return publicUser, true, nil
}

func GetUserPoints(Uid string) (int, bool, error) {
	rows, err := DB.Query("SELECT Points FROM UserTable WHERE UserID = ?", Uid)
	points := 0
	if err != nil {
		log.Printf("GetUserPoints() #1: %v", err)
		rows.Close()
		return points, false, err
	}
	counter := 0
	for rows.Next() {
		counter += 1
		rows.Scan(&points)
	}
	rows.Close()

	return points, counter == 1, err

}

// Add user into DB
func AddUser(u User) (bool, error) {
	socialCode, err := generateSocialCode()
	if err != nil {
		log.Printf("AddUser() #1: %v", err)
		return false, err
	}

	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("AddUser() #2: %v", err)
		return false, err
	}
	defer tx.Rollback() // abort transaction if error

	stmt, err := tx.Preparex("INSERT INTO UserTable (UserID, Username, Picture, Points, Bossid, SocialCode) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("AddUser() #3: %v", err)
		return false, err
	}
	defer stmt.Close() //defer the closing of SQL statement to ensure it Closes once the function completes

	_, err = stmt.Exec(u.UserID, u.Username, u.Picture, u.Points, u.BossId, socialCode)
	if err != nil {
		log.Printf("AddUser(): breaky 4: %v", err)
		return false, err
	}

	tx.Commit() //commit transaction to database

	return true, nil
}

// Generates a public code to differentiate users
func generateSocialCode() (string, error) {
	codearr := make([]byte, USER_CODE_LEN)

	// Seed at the current time
	randgen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Loop until a unique ID was created
	isUnique := false
	var code string
	for !isUnique {
		// Generate a code
		for i := range codearr {
			codearr[i] = socialcode_set[randgen.Intn(len(socialcode_set))]
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
		SET Username = ?, Points = ?, Bossid = ?
		WHERE UserID = ?
	`)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.Username, u.Points, u.BossId, uid)
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

// Search for one specific user by their social code.
func SearchUserCode(code string, includeID bool) (map[string]interface{}, bool, error) {
	rows, err := DB.Query("SELECT UserID FROM UserTable WHERE SocialCode=?;", code)
	var user map[string]interface{}
	var uid string
	var found bool = false
	if err != nil {
		log.Printf("SearchUserCode() #1: %v", err)
		return user, false, err
	}

	counter := 0
	for rows.Next() {
		counter += 1

		err := rows.Scan(&uid)
		if err != nil {
			log.Printf("SearchUserCode() #2: %v", err)
			return map[string]interface{}{}, false, err
		}
	}
	rows.Close()

	user, found, err = GetPublicUser(uid)
	if !found || err != nil {
		log.Printf("SearchUserCode() #4: did not find a user: %v", err)
		return user, false, err
	}

	if includeID {
		user["UserID"] = uid
	}

	return user, counter == 1, nil
}

// Search for any users that match this username.
func SearchUsername(uname string, includeID bool) ([]map[string]interface{}, bool, error) {
	rows, err := DB.Query("SELECT UserID FROM UserTable WHERE Username LIKE ?", "%"+uname+"%")
	var users []map[string]interface{}
	var userids []string
	if err != nil {
		log.Printf("SearchUsername() #1: %v", err)
		return users, false, err
	}

	counter := 0
	for rows.Next() {
		counter += 1

		var uid string
		err := rows.Scan(&uid)
		if err != nil {
			log.Printf("SearchUsername() #2: %v", err)
			return users, false, err
		}

		userids = append(userids, uid)
	}
	rows.Close()

	for _, uID := range userids {
		user, found, err := GetPublicUser(uID)
		if !found || err != nil {
			log.Printf("SearchUsername() #3: did not find a user: %v", err)
			return users, false, err
		}

		if includeID {
			user["UserID"] = uID
		}

		users = append(users, user)
	}

	// Return if found any matches
	return users, counter > 0, nil
}

func AddFriend(my_uid string, their_soccode string) (bool, error) {
	their_user, found, err := SearchUserCode(their_soccode, true)
	if !found || err != nil {
		log.Printf("AddFriend() #1: could not find other user: %v", err)
		return false, err
	}
	their_uid, ok := their_user["UserID"].(string)
	if !ok {
		log.Printf("AddFriend(): conversion error")
		return false, err
	}

	tx, err := DB.Beginx()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() // aborrt transaction if error

	stmt, err := tx.Preparex("INSERT INTO Friends (userA, userB) VALUES (?, ?)")
	if err != nil {
		log.Printf("AddFriend() #2: error preparing statement: %v", err)
		return false, err
	}
	defer stmt.Close()

	// Order by string compare to avoid duplicate rows
	var firstID, secondID string
	if my_uid < their_uid {
		firstID, secondID = my_uid, their_uid
	} else {
		firstID, secondID = their_uid, my_uid
	}

	_, err = stmt.Exec(firstID, secondID)
	if err != nil {
		log.Printf("AddFriend() #3: error adding friend pair: %v", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteFriend(my_uid string, their_soccode string) (bool, error) {
	their_user, found, err := SearchUserCode(their_soccode, true)
	if !found || err != nil {
		log.Printf("DeleteFriend() #1: could not find other user: %v", err)
		return false, err
	}
	their_uid, ok := their_user["UserID"].(string)
	if !ok {
		log.Printf("AddFriend(): conversion error")
		return false, err
	}

	// Order enforced in schema
	var firstID, secondID string
	if my_uid < their_uid {
		firstID, secondID = my_uid, their_uid
	} else {
		firstID, secondID = their_uid, my_uid
	}

	tx, err := DB.Beginx()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() // aborrt transaction if error

	// Depends which user was denoted as userA vs. userB
	stmt, err := tx.Preparex("DELETE FROM Friends WHERE ? = userA AND ? = userB")
	if err != nil {
		log.Printf("DeleteFriend() #2: error preparing statement: %v", err)
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(firstID, secondID)
	if err != nil {
		log.Printf("DeleteFriend() #3: error adding friend pair: %v", err)
		return false, err
	}

	tx.Commit()

	return true, nil
}

// Returns all the friends of the current user
func GetFriendList(my_uid string, includeID bool) ([]map[string]interface{}, error) {
	friends := []map[string]interface{}{}
	friendIDs := []string{}

	rows, err := DB.Query("SELECT * FROM Friends WHERE userA=? OR userB=?;", my_uid, my_uid)
	if err != nil {
		log.Printf("GetFriendList() #1: %v", err)
		return friends, err
	}

	counter := 0
	for rows.Next() {
		counter += 1

		var userAid, userBid string
		err := rows.Scan(&userAid, &userBid)
		if err != nil {
			log.Printf("GetFriendList() #2: %v", err)
			return friends, err
		}

		if my_uid == userAid {
			friendIDs = append(friendIDs, userBid)
		} else {
			friendIDs = append(friendIDs, userAid)
		}
	}
	rows.Close()

	for _, fID := range friendIDs {
		fUser, found, err := GetPublicUser(fID)
		if !found || err != nil {
			log.Printf("GetFriendList(): could not retrieve friend: %v", err)
			return friends, err
		}

		if includeID {
			fUser["UserID"] = fID
		}

		friends = append(friends, fUser)
	}

	return friends, nil
}
