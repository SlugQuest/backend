package unit_tests

import (
	"testing"

	. "slugquest.com/backend/crud"
)


var testTeamUser = User{
	UserID: "not_a_real_user",
	Username: "asdf", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: 0 ,    
	BossId: 0,
	SocialCode: "testTeamUser0", // Friendly code to uniquely identify (public)
}



var testTeamUser2 = User{
	UserID: "not_a_real_user2",
	Username: "xkcd", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: 0,    
	BossId: 0,
	SocialCode: "testTeamUser1", // Friendly code to uniquely identify (public)
}
var testTeamUser3 = User{
	UserID: "not_a_real_user3",
	Username: "asdf", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: 0 ,    
	BossId: 0,
	SocialCode: "testTeamUser3", // Friendly code to uniquely identify (public)
}

var testTeamUser4 = User{
	UserID: "not_a_real_user4",
	Username: "asdf", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: 0 ,    
	BossId: 0,
	SocialCode: "testTeamUser4", // Friendly code to uniquely identify (public)
}

var testTeamUser5 = User{
	UserID: "not_a_real_user5",
	Username: "asdf", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: 0 ,    
	BossId: 0,
	SocialCode: "testTeamUser5", // Friendly code to uniquely identify (public)
}
var testTeamUser6 = User{
	UserID: "not_a_real_user5",
	Username: "asdf", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: 0 ,    
	BossId: 0,
	SocialCode: "testTeamUser6", // Friendly code to uniquely identify (public)
}
var testTeamUser7 = User{
	UserID: "not_a_real_user7",
	Username: "asdf", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: 0 ,    
	BossId: 0,
	SocialCode: "testTeamUser7", // Friendly code to uniquely identify (public)
}

var testTeamUser8 = User{
	UserID: "not_a_real_user8",
	Username: "asdf", // Set by user, can be exposed
	Picture: "url",  // A0 stores their profile pics as URLs
	Points: 0 ,    
	BossId: 0,
	SocialCode: "testTeamUser8", // Friendly code to uniquely identify (public)
}


func testTeamUserOnCreateTeam(t *testing.T) {
	AddUser(testTeamUser)
	check, teamnid, err := CreateTeam("testingTeam", testTeamUser.UserID)
	if( err != nil || !check){
		t.Error("testTeamUserOnCreateTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("testTeamUserOnCreateTeam(): failed create less than 0")
	}
	user, err2 := GetTeamUsers(teamnid)

	if( err2 != nil){
		t.Error("GetUserTeam(): getting user team failed", err2)
	}
	if (len(user)!= 1){
		t.Error("GetTeamUsers(): failed asgn user not 1")
	}
}



func TestGetUserTeams(t *testing.T){
	AddUser(testTeamUser2)
	check, teamnid, err := CreateTeam("testingTeam", testTeamUser2.UserID)
	if( err != nil || !check){
		t.Error("testTeamUserOnCreateTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("testTeamUserOnCreateTeam(): failed create less than 0")
	}

	teaml, err2 := GetUserTeams(testTeamUser2.UserID)
	if (err2 != nil ){
		t.Error("GetUserTeams(): failed to GetUserTeams", err2)
	}
	if (len(teaml) != 1){
		t.Error("GetUserTeams(): failed user Team Count")
	}

}

func RemoveUserFromTeamTest(t *testing.T){
	AddUser(testTeamUser3)
	check, teamnid, err := CreateTeam("testingTeam", testTeamUser3.UserID)
	if( err != nil || !check){
		t.Error("RemoveUSErfRomCreateFail(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("RemoveUserFromCreateFail(): failed create less than 0")
	}

	boo, err2 := RemoveUserFromTeam(teamnid, testTeamUser3.SocialCode)
	if(!boo || err2 != nil){
		t.Error("RemoveUserFromTeam(): broken", err2)
	}
	user, err3 := GetTeamUsers(teamnid)

	if( err3 != nil){
		t.Error("GetUserTeam(): getting user team failed", err3)
	}
	if (len(user)!= 1){
		t.Error("GetTeamUsers(): failed asgn user not 1")
	}
}
func DeleteTeamTest(t *testing.T){
	AddUser(testTeamUser4)
	check, teamnid, err := CreateTeam("testingTeam", testTeamUser4.UserID)
	if( err != nil || !check){
		t.Error("DeleteTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("DeleteTeam(): failed create less than 0")
	}

	boo, err2 := DeleteTeam(teamnid)
	if(!boo || err2 != nil){
		t.Error("RemoveUserFromTeam(): broken", err2)
	}

	teaml, err2 := GetUserTeams(testTeamUser4.UserID)
	if (err2 != nil ){
		t.Error("DeleteTeam(): failed to GetUserTeams", err2)
	}
	if (len(teaml) != 0){
		t.Error("DeleteTeam(): failed user Team Count")
	}
}

func testTeamUserAddTeam(t *testing.T){
	AddUser(testTeamUser5)
	AddUser(testTeamUser6)
	check, teamnid, err := CreateTeam("testingTeam", testTeamUser5.UserID)
	if( err != nil || !check){
		t.Error("testTeamUserAddTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("testTeamUserAddTeamFail(): failed create less than 0")
	}
	boo, err1 := AddUserToTeam(teamnid, testTeamUser6.SocialCode);
	if (!boo || err1 != nil){
		t.Error("AddUserTeam(): adderr", err1)
	}
	teaml, err3 := GetUserTeams(testTeamUser6.UserID)
	if(err3 != nil){
		t.Error("TestAddUserTeam(): query fail ", err3)
	}
	if(len(teaml) != 0){
		t.Error("TestAddUserTeam(): not added correctly")
	}

}


func testTeamUserAddTeamUid(t *testing.T){
	AddUser(testTeamUser7)
	AddUser(testTeamUser2)
	check, teamnid, err := CreateTeam("testingTeam", testTeamUser7.UserID)
	if( err != nil || !check){
		t.Error("testTeamUserAddTeam(): failed create error", err)
	}
	if(teamnid < 0){
		t.Error("testTeamUserAddTeamFail(): failed create less than 0")
	}
	boo, err1 := AddUserToTeamUid(teamnid, testTeamUser8.UserID);
	if (!boo || err1 != nil){
		t.Error("AddUserTeamUid(): adderr", err1)
	}
	teaml, err3 := GetUserTeams(testTeamUser8.UserID)
	if(err3 != nil){
		t.Error("TestAddUserTeamUid(): query fail ", err3)
	}
	if(len(teaml) != 0){
		t.Error("TestAddUserTeamUid(): not added correctly")
	}
}