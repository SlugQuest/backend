package unit_tests

import (
	"testing"
	"time"

	. "slugquest.com/backend/crud"
)


var testUser = User{
	UserID: "not_a_real_user",
	Username: "asdf", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: "0" ,    
	BossId: "0",
	SocialCode: "testuser0" // Friendly code to uniquely identify (public)
}


var testUser2 = User{
	UserID: "not_a_real_user2",
	Username: "xkcd", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: "0" ,    
	BossId: "0",
	SocialCode: "testuser1" // Friendly code to uniquely identify (public)
}


func LoadUserIfNotExist(userForUserTable User) (bool, error) {
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

func TestUserOnCreateTeam(t *testing) {
	LoadUserIfNotExist(testUser)
	check, teamnid, err := CreateTeam("testingTeam", testUser.UserID)
	if( err != nil || !check){
		t.Error("TestUserOnCreateTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("TestUserOnCreateTeam(): failed create less than 0")
	}
	user, err2 := GetTeamUsers(teamnid)

	if( err2 != nil){
		t.Error("GetUserTeam(): getting user team failed", err2)
	}
	if (len(user)!= 1){
		t.Error("GetTeamUsers(): failed asgn user not 1")
	}
}



func TestGetUserTeams(t *testing){
	LoadUserIfNotExist(testUser)
	check, teamnid, err := CreateTeam("testingTeam", testUser.UserID)
	if( err != nil || !check){
		t.Error("TestUserOnCreateTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("TestUserOnCreateTeam(): failed create less than 0")
	}

	teaml, err2 = GetUserTeams(testUser.UserID)
	if (err2 != nil ){
		t.Error("GetUserTeams(): failed to GetUserTeams", err2)
	}
	if (teaml.len() ){
		t.Error("GetUserTeams(): failed user Team Count")
	}

}

func RemoveUserFromTeam(t *testing){
	LoadUserIfNotExist(testUser)
	check, teamnid, err := CreateTeam("testingTeam", testUser.UserID)
	if( err != nil || !check){
		t.Error("RemoveUSErfRomCreateFail(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("RemoveUserFromCreateFail(): failed create less than 0")
	}

	boo, err2; := RemoveUserFromTeam(teamnid, testUser.SocialCode)
	if(!boo || err2 != nil){
		t.Error("RemoveUserFromTeam(): broken", err2)
	}
	user, err2 := GetTeamUsers(teamnid)

	if( err3 != nil){
		t.Error("GetUserTeam(): getting user team failed", err3)
	}
	if (len(user)!= 1){
		t.Error("GetTeamUsers(): failed asgn user not 1")
	}
}
func DeleteTeamTest(t *testing){
	LoadUserIfNotExist(testUser)
	check, teamnid, err := CreateTeam("testingTeam", testUser.UserID)
	if( err != nil || !check){
		t.Error("DeleteTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("DeleteTeam(): failed create less than 0")
	}

	boo, err2; := DeleteTeam(teamnid)
	if(!boo || err2 != nil){
		t.Error("RemoveUserFromTeam(): broken", err2)
	}

	teaml, err2 = GetUserTeams(testUser.UserID)
	if (err2 != nil ){
		t.Error("DeleteTeam(): failed to GetUserTeams", err2)
	}
	if (teaml.len() ){
		t.Error("DeleteTeam(): failed user Team Count")
	}
}

func TestUserAddTeam(t *testing){
	LoadUserIfNotExist(testUser)
	LoadUserIfNotExist(testUser2)
	check, teamnid, err := CreateTeam("testingTeam", testUser.UserID)
	if( err != nil || !check){
		t.Error("TestUserAddTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("TestUserAddTeamFail(): failed create less than 0")
	}
	boo, err1 := AddUserToTeam(teamnid, testUser2.SocialCode);
	if (!boo || err1 != nil){
		t.Error("AddUserTeam(): adderr", err1)
	}
	teaml, err3 = GetUserTeams(testUser2.UserID)
	if(err3 != nil){
		t.Error("TestAddUserTeam(): query fail ", err3)
	}
	if(teaml.len() != 0){
		t.Error("TestAddUserTeam(): not added correctly")
	}

}


func TestUserAddTeamUid(t *testing){
	LoadUserIfNotExist(testUser)
	LoadUserIfNotExist(testUser2)
	check, teamnid, err := CreateTeam("testingTeam", testUser.UserID)
	if( err != nil || !check){
		t.Error("TestUserAddTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("TestUserAddTeamFail(): failed create less than 0")
	}
	boo, err1 := AddUserToTeamUid(teamnid, testUser2.UserID);
	if (!boo || err1 != nil){
		t.Error("AddUserTeamUid(): adderr", err1)
	}
	teaml, err3 = GetUserTeams(testUser2.UserID)
	if(err3 != nil){
		t.Error("TestAddUserTeamUid(): query fail ", err3)
	}
	if(teaml.len() != 0){
		t.Error("TestAddUserTeamUid(): not added correctly")
	}
}