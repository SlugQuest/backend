package unit_tests

import (
	"log"
	"testing"

	. "slugquest.com/backend/crud"
)

var testTeamUser1 = User{
	UserID:   "not_a_real_user1",
	Username: "asdf",
	Picture:  "url",
	Points:   0,
	BossId:   1,
}

var testTeamUser2 = User{
	UserID:   "not_a_real_user2",
	Username: "xkcd",
	Picture:  "url",
	Points:   0,
	BossId:   1,
}
var testTeamUser3 = User{
	UserID:   "not_a_real_user3",
	Username: "asdf",
	Picture:  "url",
	Points:   0,
	BossId:   1,
}

var testTeamUser4 = User{
	UserID:   "not_a_real_user4",
	Username: "asdf",
	Picture:  "url",
	Points:   0,
	BossId:   1,
}

var testTeamUser5 = User{
	UserID:   "not_a_real_user5",
	Username: "asdf",
	Picture:  "url",
	Points:   0,
	BossId:   1,
}
var testTeamUser6 = User{
	UserID:   "not_a_real_user6",
	Username: "asdf",
	Picture:  "url",
	Points:   0,
	BossId:   1,
}
var testTeamUser7 = User{
	UserID:   "not_a_real_user7",
	Username: "asdf",
	Picture:  "url",
	Points:   0,
	BossId:   1,
}

var testTeamUser8 = User{
	UserID:   "not_a_real_user8",
	Username: "asdf",
	Picture:  "url",
	Points:   0,
	BossId:   1,
}

func addTestTeamUsers() (bool, error) {
	users := []*User{
		&testTeamUser1,
		&testTeamUser2,
		&testTeamUser3,
		&testTeamUser4,
		&testTeamUser5,
		&testTeamUser6,
		&testTeamUser7,
		&testTeamUser8,
	}

	for _, u := range users {
		fromDB, found, getErr := GetUser(u.UserID)
		if getErr != nil {
			log.Printf("Could not find test users for team unit tests: %v", getErr)
			return false, getErr
		}

		// Set social code in their struct
		if found {
			u.SocialCode = fromDB.SocialCode
		} else {
			added, err := AddUser(*u)
			if !added || err != nil {
				log.Printf("Could not add test users for team unit tests: %v", err)
				return false, err
			}

			fromDB, found, getErr := GetUser(u.UserID)
			if !found || getErr != nil {
				log.Printf("Could not find test users for team unit tests: %v", getErr)
				return false, getErr
			}

			u.SocialCode = fromDB.SocialCode
		}
	}

	return true, nil
}

func TestTeamUserOnCreateTeam(t *testing.T) {
	setupSuccess, err := addTestTeamUsers()
	if !setupSuccess || err != nil {
		t.Error("Could not setup test suite", err)
	}

	check, teamnid, err := CreateTeam("testingTeam", testTeamUser1.UserID)
	if err != nil || !check {
		t.Error("TestTeamUserOnCreateTeam(): failed create error", err)
	}
	if teamnid < 0 {
		t.Error("TestTeamUserOnCreateTeam(): failed create less than 0")
	}
	user, err2 := GetTeamUsers(teamnid)

	if err2 != nil {
		t.Error("TestTeamUserOnCreateTeam(): getting user team failed", err2)
	}
	if len(user) != 1 {
		t.Error("TestTeamUserOnCreateTeam(): failed asgn user not 1")
	}
}

func TestGetUserTeams(t *testing.T) {
	setupSuccess, err := addTestTeamUsers()
	if !setupSuccess || err != nil {
		t.Error("Could not setup test suite", err)
	}

	check, teamnid, err := CreateTeam("testingTeam", testTeamUser2.UserID)
	if err != nil || !check {
		t.Error("TestGetUserTeams(): failed create error", err)
	}
	if teamnid < 0 {
		t.Error("TestGetUserTeams(): failed create less than 0")
	}

	teaml, err2 := GetUserTeams(testTeamUser2.UserID)
	if err2 != nil {
		t.Error("TestGetUserTeams(): failed to GetUserTeams", err2)
	}
	if len(teaml) != 1 {
		t.Error("TestGetUserTeams(): failed user Team Count")
	}

}

func TestRemoveUserFromTeam(t *testing.T) {
	setupSuccess, err := addTestTeamUsers()
	if !setupSuccess || err != nil {
		t.Error("Could not setup test suite", err)
	}

	check, teamnid, err := CreateTeam("testingTeam", testTeamUser3.UserID)
	if err != nil || !check {
		t.Error("TestRemoveUserFromTeam(): failed create error", err)
	}
	if teamnid < 0 {
		t.Error("TestRemoveUserFromTeam(): failed create less than 0")
	}

	boo, err2 := RemoveUserFromTeam(teamnid, testTeamUser3.SocialCode)
	if !boo || err2 != nil {
		t.Error("TestRemoveUserFromTeam(): broken", err2)
	}
	usersinteam, err3 := GetTeamUsers(teamnid)

	if err3 != nil {
		t.Error("TestRemoveUserFromTeam(): getting user team failed", err3)
	}
	if len(usersinteam) != 0 {
		t.Errorf("TestRemoveUserFromTeam(): team should now be empty, got %v", len(usersinteam))
	}
}

func TestDeleteTeam(t *testing.T) {
	setupSuccess, err := addTestTeamUsers()
	if !setupSuccess || err != nil {
		t.Error("Could not setup test suite", err)
	}

	check, teamnid, err := CreateTeam("testingTeam", testTeamUser4.UserID)
	if err != nil || !check {
		t.Error("DeleteTeam(): failed create error", err)
	}
	if teamnid < 0 {
		t.Error("DeleteTeam(): failed create less than 0")
	}

	boo, err2 := DeleteTeam(teamnid)
	if !boo || err2 != nil {
		t.Error("RemoveUserFromTeam(): broken", err2)
	}

	teaml, err2 := GetUserTeams(testTeamUser4.UserID)
	if err2 != nil {
		t.Error("DeleteTeam(): failed to GetUserTeams", err2)
	}
	if len(teaml) != 0 {
		t.Errorf("DeleteTeam(): failed team user count, got %v instead of 0", len(teaml))
	}
}

func TestTeamUserAddTeam(t *testing.T) {
	setupSuccess, err := addTestTeamUsers()
	if !setupSuccess || err != nil {
		t.Error("Could not setup test suite", err)
	}

	check, teamnid, err := CreateTeam("testingTeam", testTeamUser5.UserID)
	if err != nil || !check {
		t.Error("TestTeamUserAddTeam(): failed create error", err)
	}
	if teamnid < 0 {
		t.Error("TestTeamUserAddTeam(): failed create less than 0")
	}
	boo, err1 := AddUserToTeam(teamnid, testTeamUser6.SocialCode)
	if !boo || err1 != nil {
		t.Error("TestTeamUserAddTeam(): erorr adding user to team", err1)
	}
	teaml, err3 := GetUserTeams(testTeamUser6.UserID)
	if err3 != nil {
		t.Error("TestAddUserTeam(): query fail ", err3)
	}

	if len(teaml) != 1 {
		t.Errorf("TestAddUserTeam(): team should have 1 user, got %v", len(teaml))
	}
}

func TestTeamUserAddTeamUid(t *testing.T) {
	setupSuccess, err := addTestTeamUsers()
	if !setupSuccess || err != nil {
		t.Error("Could not setup test suite", err)
	}

	check, teamnid, err := CreateTeam("testingTeam", testTeamUser7.UserID)
	if err != nil || !check {
		t.Error("TestTeamUserAddTeamUid(): failed create error", err)
	}

	if teamnid < 0 {
		t.Error("TestTeamUserAddTeamUid(): failed create less than 0")
	}

	boo, err1 := AddUserToTeamUid(teamnid, testTeamUser8.UserID)
	if !boo || err1 != nil {
		t.Error("TestTeamUserAddTeamUid(): adderr", err1)
	}

	teaml, err3 := GetUserTeams(testTeamUser8.UserID)
	if err3 != nil {
		t.Error("TestTeamUserAddTeamUid(): query fail ", err3)
	}

	if len(teaml) != 1 {
		t.Errorf("TestTeamUserAddTeamUid(): team should have 1 user, got %v", len(teaml))
	}
}
