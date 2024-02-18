package testing

import (
	"log"

	. "slugquest.com/backend/crud"
)

var userForUserTable = User{
	UserID:   "usertable_user_id",
	Username: "sluggo2",
	Picture:  "lmao.jpg",
	Points:   1,
	BossId:   1,
}

func TestUPoints() bool {
	// NEEDS TO BE DONE
	return false
}

func TestAddUser() bool {
	addSuccess, addErr := AddUser(userForUserTable)
	if addErr != nil || !addSuccess {
		log.Printf("TestAddUser(): couldn't add user")
		return false
	}

	_, found, _ := GetUser(userForUserTable.UserID)
	if !found {
		log.Println("TestAddUser(): add failed")
		return false
	}

	return true
}

func TestEditUser() bool {
	editedUser := User{
		UserID:   userForUserTable.UserID,
		Username: "not in DB, not tested",
		Picture:  "not in DB, not tested",
		Points:   5,
		BossId:   10,
	}

	editSuccess, editErr := EditUser(editedUser, editedUser.UserID)
	if editErr != nil || !editSuccess {
		log.Printf("TestEditUser(): error editing user: %v", editErr)
		return false
	}

	checkE, _, _ := GetUser(editedUser.UserID)
	if checkE.Points != 5 || checkE.BossId != 10 {
		log.Println("TestEditUser(): edit verfication failed")
		return false
	}

	return true
}

func TestDeleteUser() bool {
	deleteSuccess, deleteErr := DeleteUser(userForUserTable.UserID)
	if deleteErr != nil || !deleteSuccess {
		log.Printf("TestDeleteUser(): couldn't delete user")
		return false
	}

	_, found, _ := GetUser(userForUserTable.UserID)
	if found {
		log.Println("TestDeleteUser(): delete failed")
		return false
	}

	return true
}
