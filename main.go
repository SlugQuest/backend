package main

import (
	"log"

	envfuncs "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"slugquest.com/backend/authentication"
	"slugquest.com/backend/crud"
	_"slugquest.com/backend/testing"
)

func main() {
	// Load .env
	if env_err := envfuncs.Load(); env_err != nil {
		log.Fatalf("main(): Error loading the .env file: %v", env_err)
	}

	// Create new authenticator to pass to the router
	auth, auth_err := authentication.NewAuthenticator()
	if auth_err != nil {
		log.Fatalf("main(): Failed to initialize the authenticator: %v", auth_err)
	}
	router := CreateRouter(auth)
	// utest := testing.RunAllTests()
	// if !utest {
	// 	log.Println("main(): unit test failure")
	// 	return
	// }


	conn_err := crud.ConnectToDB(false)
	if conn_err != nil {
		log.Fatalf("main(): Error connecting to database: %v", conn_err)
		return
	}


	log.Print("Running at http://localhost:8080")
	router_err := router.Run() // listen and serve on 0.0.0.0:8080
	if router_err != nil {
		log.Fatalf("main(): couldn't run server gg: %v", router_err)
	}
}
