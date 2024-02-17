package crud

import "fmt"

// GetBossById retrieves boss information by BossID.
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

func GetCurrBossHealth(uid string) (int, error) {
	user, exists, err := GetUser(uid)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, fmt.Errorf("User not found")
	}

	boss, exists, err := GetBossById(user.BossId)
	if err != nil {
		return 0, err
	}

	if !exists {
		fmt.Println("Naur")
		return 0, fmt.Errorf("no boss found")
	}

	currBossHealth := boss.Health - user.Points
	fmt.Printf("in crud: currBossHealth: %v\n", currBossHealth)

	if currBossHealth < 0 { //should never get here, pass task has logic to update boss id
		currBossHealth = 0
	}

	return currBossHealth, nil
}