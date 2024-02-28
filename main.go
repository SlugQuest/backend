package main

import (
	_"fmt"
	"log"
	_"time"

	envfuncs "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"slugquest.com/backend/authentication"
	"slugquest.com/backend/crud"
)

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

	// go func() { //launching a goroutine
	// 	timer := time.NewTimer(0) // Initial trigger
	// 	for {                     //loop that runs forever

	// 		// Block until timer finishes. When done, it sends a message on the channel
	// 		// timer.C; no other code in this goroutine is executed until that happens.
	// 		<-timer.C
	// 		fmt.Println("Running task population...")

	// 		err := crud.PopRecurringTasksMonth()
	// 		if err != nil {
	// 			fmt.Printf("Error populating recurring tasks: %v\n", err)
	// 		}

	// 		// Reschedule for the next month
	// 		nextMonth := time.Now().AddDate(0, 1, 0)
	// 		firstDayOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	// 		timeUntilNextMonth := time.Until(firstDayOfMonth)
	// 		timer.Reset(timeUntilNextMonth)
	// 		//Stops a ticker and resets its period to the specified duration.
	// 		//The next tick will arrive after the new period elapses.
	// 	}
	// }()

	log.Print("Running at http://localhost:8080")
	router_err := router.Run() // listen and serve on 0.0.0.0:8080
	if router_err != nil {
		log.Fatalf("main(): couldn't run server gg: %v", router_err)
	}
}
