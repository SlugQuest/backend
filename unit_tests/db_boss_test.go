package unit_tests

import (
	"path/filepath"
	"testing"

	. "slugquest.com/backend/crud"
)

var testBoss = Boss{
	BossID: testUser.BossId,
	Name:   "testboss_name",
	Health: 30,
	Image:  filepath.Join("images", "clown.jpeg"),
}

var secondTestBoss = Boss{
	BossID: testUser.BossId + 1,
	Name:   "secondtestboss_name",
	Health: 30,
	Image:  filepath.Join("images", "lol.jpeg"),
}

func TestGetAddBoss(t *testing.T) {
	// Added as dummy data, should be found
	returnedBoss, found, getBossErr := GetBossById(testBoss.BossID)
	if getBossErr != nil {
		t.Errorf("TestGetAddBoss(): error querying for boss: %v", getBossErr)
	}

	if !found {
		t.Error("TestGetAddBoss(): error finding boss")
	}

	if returnedBoss.BossID != testBoss.BossID {
		t.Error("TestGetAddBoss(): found wrong boss")
	}
}

func TestAddBoss(t *testing.T) {
	addBossSuccess, addBossErr := AddBoss(secondTestBoss)
	if addBossErr != nil || !addBossSuccess {
		t.Errorf("TestGetAddBoss(): error adding test boss: %v", addBossErr)
	}
}

func TestGetBossId(t *testing.T) {
	boss, found, err := GetBossById(testBoss.BossID)

	if err != nil {
		t.Errorf("TestGetBossId(): error getting test boss: %v", err)
	}

	if !found {
		t.Error("TestGetBossId(): didn't find boss")
	}

	if boss.BossID != testBoss.BossID {
		t.Errorf("TestGetBossId(): found wrong boss, expected %v, got %v", testBoss.BossID, boss.BossID)
	}
}

func TestGetCurrBossHealth(t *testing.T) {
	currBossHealth, err := GetCurrBossHealth(testUser.UserID)
	if err != nil {
		t.Errorf("TestGetCurrBossHealth(): error getting current boss health: %v", err)
	}

	expectedHealth := testBoss.Health - testUser.Points
	if currBossHealth != expectedHealth {
		t.Errorf("TestGetCurrBossHealth(): returned wrong health, expected %v, got %v", expectedHealth, currBossHealth)
	}
}
