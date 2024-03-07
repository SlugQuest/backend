package crud

import (
	"fmt"
	"log"
)

// AddBoss adds a new boss to the BossTable.
// Inputs:
// boss - a Boss struct representing the boss to be added
// Outputs:
// bool  - a success flag indicating whether the boss addition was successful
// error - any error that occurred during the transaction or statement execution
func AddBoss(boss Boss) (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("AddBoss(): error beginning transaction")
		return false, err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(`
		INSERT INTO BossTable (BossID, BossName, Health, BossImage)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		log.Printf("AddBoss(): error adding boss to table")
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(boss.BossID, boss.Name, boss.Health, boss.Image)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("AddBoss(): error committing transaction")
		return false, err
	}

	return true, nil
}

// PopBossTable checks if the BossTable is populated and adds default bosses if not.
// Inputs: None
// Outputs:
//
//	bool  - a success flag indicating whether the BossTable population was successful
//	error - any error that occurred during the transaction or statement execution
func PopBossTable() (bool, error) {
	tx, err := DB.Beginx()
	if err != nil {
		log.Printf("popBossTable(): error beginning transaction")
		return false, err
	}
	defer tx.Rollback()
	// Checking if BossTable is already populated
	var count int
	err = tx.Get(&count, "SELECT COUNT(*) FROM BossTable")
	if err != nil {
		log.Printf("popBossTable(): error checking BossTable population, %v ", err)
		return false, err
	}

	tx.Commit()

	if count == 0 {
		// Add default bosses if BossTable is not populated
		basicBosses := []Boss{
			{BossID: 1, Name: "Default Boss 1", Health: 40, Image: "clown.jpeg"},
			{BossID: 2, Name: "Default Boss 2", Health: 90, Image: "clown.jpeg"},
			{BossID: 3, Name: "Default Boss 3", Health: 200, Image: "clown.jpeg"},
			{BossID: 4, Name: "Default Boss 4", Health: 300, Image: "clown.jpeg"},
		}

		for _, boss := range basicBosses {
			_, err := AddBoss(boss)
			if err != nil {
				log.Printf("PopBossTable(): error adding boss to BossTable: %v", err)
				return false, err
			}
		}
	}

	return true, nil
}

// GetBossById retrieves boss information by BossID.
// Inputs:
// bossID - an integer representing the BossID
// Outputs:
// Boss - the retrieved boss information
// bool - a success flag indicating whether the boss was found
// error - any error that occurred during the query
func GetBossById(bossID int) (Boss, bool, error) {
	var boss Boss
	rows, err := DB.Query("SELECT * FROM BossTable WHERE BossID = ?", bossID)
	if err != nil {
		return boss, false, err
	}

	counter := 0
	for rows.Next() {
		counter += 1
		if err := rows.Scan(&boss.BossID, &boss.Name, &boss.Health, &boss.Image); err != nil {
			return boss, false, err
		}
	}

	rows.Close()

	return boss, counter == 1, nil
}

// GetCurrBossHealth calculates the current boss health based on user points.
// Inputs:
// uid - a string representing the UserID
// Outputs:
// int  - the calculated current boss health
// error - any error that occurred during the user or boss retrieval
func GetCurrBossHealth(uid string) (int, error) {
	user, exists, err := GetUser(uid)
	if err != nil {
		log.Printf("GetCurrBossHealth() #1: %v", err)
		return 0, err
	}

	if !exists {
		log.Print("GetCurrBossHealth() #2: User not found")
		return 0, fmt.Errorf("User not found")
	}

	boss, exists, err := GetBossById(user.BossId)
	if err != nil {
		log.Printf("GetCurrBossHealth() #3: %v", err)
		return 0, err
	}

	if !exists {
		log.Print("GetCurrBossHealth() #4: No boss found")
		return 0, fmt.Errorf("no boss found")
	}

	currBossHealth := boss.Health - user.Points

	if currBossHealth < 0 { //should never get here, pass task has logic to update boss id
		currBossHealth = 0
	}

	return currBossHealth, nil
}
