package unit_tests

import (
	"strconv"
	"strings"
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

func TestAddUser(t *testing.T) {
	addSuccess, addErr := AddUser(userForUserTable)
	if addErr != nil || !addSuccess {
		t.Errorf("TestAddGetUser(): couldn't add user: %v", addErr)
	}

	foundUser, found, getErr := GetUser(userForUserTable.UserID)
	if !found {
		t.Errorf("TestAddGetUser(): could not find user after adding: %v", getErr)
	}

	if userForUserTable.UserID != foundUser.UserID || userForUserTable.Points != foundUser.Points || userForUserTable.BossId != foundUser.BossId {
		t.Error("TestAddGetUser(): found wrong user")
	}
}

// To keep unit tests independent, always have the test user added
func checkIfTestUserAdded() (bool, error) {
	_, found, getErr := GetUser(userForUserTable.UserID)
	if getErr != nil {
		return false, getErr
	}

	if !found {
		addSuccess, addErr := AddUser(userForUserTable)
		if addErr != nil || !addSuccess {
			return false, addErr
		}

	}

	return true, nil
}

func TestGetUserPoints(t *testing.T) {
	checkIfTestUserAdded()

	points, found, err := GetUserPoints(userForUserTable.UserID)

	if err != nil {
		t.Errorf("TestGetUserPoints(): %v", err)
	}
	if !found {
		t.Error("TestGetUserPoints(): couldn't find user")
	}
	if points != userForUserTable.Points {
		t.Errorf("TestGetUserPoints(): wrong number of points, expected %v, got %v", userForUserTable.Points, points)
	}
}

func TestEditUser(t *testing.T) {
	checkIfTestUserAdded()

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
	if checkE.Points != editedUser.Points || checkE.BossId != editedUser.BossId {
		t.Error("TestEditUser(): edit verfication failed")
	}
}

func TestDeleteUser(t *testing.T) {
	checkIfTestUserAdded()

	deleteSuccess, deleteErr := DeleteUser(userForUserTable.UserID)
	if deleteErr != nil || !deleteSuccess {
		t.Errorf("TestDeleteUser(): couldn't delete user: %v", deleteErr)
	}

	_, found, _ := GetUser(userForUserTable.UserID)
	if found {
		t.Error("TestDeleteUser(): delete failed, found user")
	}
}

func TestSearchUserCode(t *testing.T) {
	checkIfTestUserAdded()

	// Social code generated upon add to DB
	fullUserInfo, found, err := GetUser(userForUserTable.UserID)
	if !found || err != nil {
		t.Errorf("TestSearchUserCode(): couldn't find user: %v", err)
	}

	socialcode := fullUserInfo.SocialCode
	searchedUser, found, err := SearchUserCode(socialcode)
	if !found || err != nil {
		t.Errorf("TestSearchUserCode(): couldn't search by social code: %v", err)
	}

	if userForUserTable.UserID != searchedUser.UserID || userForUserTable.Points != searchedUser.Points || userForUserTable.BossId != searchedUser.BossId {
		t.Error("TestSearchUserCode(): found wrong user")
	}
}

func TestSearchUsername(t *testing.T) {
	common := "common_name"
	numUsers := 5
	for i := 1; i <= numUsers; i++ {
		user := User{
			UserID:   "user" + strconv.Itoa(i),
			Username: common + strconv.Itoa(i),
			Picture:  strconv.Itoa(i) + ".png",
			Points:   i,
			BossId:   1,
		}

		addSuccess, addErr := AddUser(user)
		if addErr != nil || !addSuccess {
			t.Errorf("TestSearchUsername(): couldn't add user: %v", addErr)
		}
	}

	foundUsers, found, err := SearchUsername(common)
	if !found || err != nil {
		t.Errorf("TestSearchUsername(): didn't find any users on search: %v", err)
	}

	if len(foundUsers) != numUsers {
		t.Errorf("TestSearchUsername(): search did not return correct num of users, expected %v, got %v", numUsers, len(foundUsers))
	}

	for _, user := range foundUsers {
		if !strings.Contains(user.Username, common) {
			t.Errorf("TestSearchUsername(): username did not contain queried string")
		}
	}
}

func TestMultipleUserLifecycle(t *testing.T) {
	// Add multiple users to ensure no constraints break
	for i := 1; i <= 10; i++ {
		user := User{
			UserID:   "adduser" + strconv.Itoa(i),
			Username: "newuser" + strconv.Itoa(i),
			Picture:  strconv.Itoa(i) + ".png",
			Points:   i,
			BossId:   1,
		}

		addSuccess, addErr := AddUser(user)
		if addErr != nil || !addSuccess {
			t.Errorf("TestMultipleUserLifecycle(): couldn't add user: %v", addErr)
		}

		foundUser, found, getErr := GetUser(user.UserID)
		if getErr != nil || !found {
			t.Errorf("TestMultipleUserLifecycle(): could not find user after adding: %v", getErr)
		}

		if user.UserID != foundUser.UserID || user.Points != foundUser.Points || user.BossId != foundUser.BossId {
			t.Error("TestMultipleUserLifecycle(): found wrong user")
		}

		deleteSuccess, deleteErr := DeleteUser(user.UserID)
		if deleteErr != nil || !deleteSuccess {
			t.Errorf("TestMultipleUserLifecycle(): couldn't delete user: %v", deleteErr)
		}

		_, found, _ = GetUser(user.UserID)
		if found {
			t.Error("TestMultipleUserLifecycle(): delete failed, found user")
		}
	}
}
