package testing

import (
	"path/filepath"
	"testing"

	. "slugquest.com/backend/crud"
)

var userForBossTable = User{
	UserID:   "bosstable_user_id",
	Username: "sluggo3",
	Picture:  "rofl.jpg",
	Points:   1,
	BossId:   1,
}

var testBoss = Boss{
	BossID: 1,
	Name:   "testboss_name",
	Health: 30,
	Image:  filepath.Join("images", "clown.jpeg"),
}

func TestAddBoss(t *testing.T) {
	addBossSuccess, addBossErr := AddBoss(testBoss)
	if addBossErr != nil || !addBossSuccess {
		t.Errorf("TestGetCurrBossHealth(): error adding test boss: %v", addBossErr)
	}
}

func TestGetCurrBossHealth(t *testing.T) {
	addUserSuccess, addUserErr := AddUser(userForBossTable)
	if addUserErr != nil || !addUserSuccess {
		t.Errorf("TestGetCurrBossHealth(): error adding test user: %v", addUserErr)
	}

	// addBossSuccess, addBossErr := AddBoss(testBoss)
	// if addBossErr != nil || !addBossSuccess {
	// 	t.Errorf("TestGetCurrBossHealth(): error adding test boss: %v", addBossErr)
	// }
	TestAddBoss(t)

	currBossHealth, err := GetCurrBossHealth(userForBossTable.UserID)
	if err != nil {
		t.Errorf("TestGetCurrBossHealth(): error getting current boss health: %v", err)
	}

	if currBossHealth != testBoss.Health {
		t.Errorf("TestGetCurrBossHealth(): returned wrong health, expected %v, got %v", testBoss.Health, currBossHealth)
	}
}
