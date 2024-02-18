package testing

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"slugquest.com/backend/crud"
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

func TestGetCurrBossHealth(t *testing.T) {
	addUserSuccess, addUserErr := crud.AddUser(userForBossTable)
	if addUserErr != nil || !addUserSuccess {
		log.Printf("TestGetCurrBossHealth(): error adding test user: %v", addUserErr)
		return false
	}

	addBossSuccess, addBossErr := AddBoss(testBoss)
	if addBossErr != nil || !addBossSuccess {
		log.Printf("TestGetCurrBossHealth(): error adding test boss: %v", addBossErr)
		return false
	}

	currBossHealth, err := crud.GetCurrBossHealth(userForBossTable.UserID)
	if err != nil {
		log.Printf("TestGetCurrBossHealth(): error getting current boss health: %v", err)
		return false
	}

	if currBossHealth != 20 {
		fmt.Printf("curr boss health: %v", currBossHealth)
		return false
	}
	return true
}

func TestAddBoss(t *testing.T) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(`
		INSERT INTO BossTable (BossID, BossName, Health, BossImage)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(boss.BossID, boss.Name, boss.Health, boss.Image)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
