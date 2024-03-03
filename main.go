package main

import (
	"fmt"
	_ "fmt"
	"log"
	"time"
	_ "time"

	envfuncs "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"slugquest.com/backend/authentication"
	"slugquest.com/backend/crud"
)

func monthlyTasks() {
	fmt.Println("Running task population...")

	err := crud.PopRecurringTasksMonth()
	if err != nil {
		fmt.Printf("Error populating recurring tasks: %v\n", err)
	}

	// Schedule the task for the next month
	nextMonth := time.Now().AddDate(0, 1, 0)
	firstDayOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	timeUntilNextMonth := time.Until(firstDayOfMonth)

	if timeUntilNextMonth < 0 {
		timeUntilNextMonth += 30 * 24 * time.Hour
	}
	time.AfterFunc(timeUntilNextMonth, monthlyTasks)

}

func main() {
	// Load .env
	if env_err := envfuncs.Load(); env_err != nil {
		log.Fatalf("main(): Error loading the .env file: %v", env_err)
		return
	}

	// Create new authenticator to pass to the router
	auth, auth_err := authentication.NewAuthenticator()
	if auth_err != nil {
		log.Fatalf("main(): Failed to initialize the authenticator: %v", auth_err)
		return
	}
	router := CreateRouter(auth)

	conn_err := crud.ConnectToDB(false)
	if conn_err != nil {
		log.Fatalf("main(): Error connecting to database: %v", conn_err)
		return
	}

	_, err := crud.PopBossTable()
	if err != nil {
		log.Fatalf("main(): Error populating bossTable: %v", conn_err)
		return
	}

	var bossCount int
	err = crud.DB.Get(&bossCount, "SELECT COUNT(*) FROM BossTable")
	if err != nil {
		log.Fatalf("main(): Error checking BossTable population: %v", err)
		return
	}

	if bossCount > 0 {
		log.Println("BossTable is populated.")
	} else {
		log.Println("BossTable is empty.")
	}

	go monthlyTasks()

	log.Print("Running at http://localhost:8080")
	router_err := router.Run() // listen and serve on 0.0.0.0:8080
	if router_err != nil {
		log.Fatalf("main(): couldn't run server gg: %v", router_err)
	}
}
