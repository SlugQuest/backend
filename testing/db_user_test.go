package testing

import (
	"testing"

	. "slugquest.com/backend/crud"
)

// Don't conflict with dummy data's user
var userForUserTable = User{
	UserID:   "usertable_user_id",
	Username: "sluggo2",
	Picture:  "lmao.jpg",
	Points:   1,
	BossId:   1,
}

// func TestUPoints(t *testing.T) {
// 	// NEEDS TO BE DONE
// }

func TestAddUser(t *testing.T) {
	addSuccess, addErr := AddUser(userForUserTable)
	if addErr != nil || !addSuccess {
		t.Errorf("TestAddUser(): couldn't add user")
	}

	_, found, _ := GetUser(userForUserTable.UserID)
	if !found {
		t.Error("TestAddUser(): add failed")
	}
}

func TestEditUser(t *testing.T) {
	editedUser := User{
		UserID:   userForUserTable.UserID,
		Username: "not in DB, not tested",
		Picture:  "not in DB, not tested",
		Points:   5,
		BossId:   10,
	}

	editSuccess, editErr := EditUser(editedUser, editedUser.UserID)
	if editErr != nil || !editSuccess {
		t.Errorf("TestEditUser(): error editing user: %v", editErr)
	}

	checkE, _, _ := GetUser(editedUser.UserID)
	if checkE.Points != 5 || checkE.BossId != 10 {
		t.Error("TestEditUser(): edit verfication failed")
	}
}

func TestDeleteUser(t *testing.T) {
	deleteSuccess, deleteErr := DeleteUser(userForUserTable.UserID)
	if deleteErr != nil || !deleteSuccess {
		t.Errorf("TestDeleteUser(): couldn't delete user")
	}

	_, found, _ := GetUser(userForUserTable.UserID)
	if found {
		t.Error("TestDeleteUser(): delete failed")
	}
}
