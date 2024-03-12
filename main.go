package main

import (
	"log"
	"os"
	"time"

	envfuncs "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"slugquest.com/backend/authentication"
	"slugquest.com/backend/crud"
)

func monthlyTasks() {
	log.Println("Running task population...")

	err := crud.PopRecurringTasksMonth()
	if err != nil {
		log.Printf("Error populating recurring tasks: %v", err)
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

	added, err := crud.AddDefaultTeam()
	if !added || err != nil {
		log.Println("Could not add default team")
	}

	// Get port to run server on, or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Print("Running at http://localhost:" + port)
	router_err := router.Run(":" + port)
	if router_err != nil {
		log.Fatalf("main(): couldn't run server gg: %v", router_err)
	}
}
