package main

import (
	"log"

	envfuncs "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"slugquest.com/backend/authentication"
	"slugquest.com/backend/crud"
	"slugquest.com/backend/testing"
)

func main() {
	// Load .env
	if env_err := envfuncs.Load(); env_err != nil {
		log.Fatalf("Error loading the .env file: %v", env_err)
		return
	}

	// Create new authenticator to pass to the router
	auth, auth_err := authentication.NewAuthenticator()
	if auth_err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", auth_err)
		return
	}
	router := CreateRouter(auth)

	// DB loading shenanigans
	conn_err := crud.ConnectToDB()
	if conn_err != nil {
		log.Fatalf("Error connecting to database: %v", conn_err)
		return
	}
	dummy_err := crud.LoadDumbData()
	if dummy_err != nil {
		log.Fatalf("error loaduing dumb data: %v", dummy_err)
		return
	}
	utest := testing.RunAllTests()
	if !utest {
		log.Fatal("unit test failure")
		return
	}

	log.Print("Running at http://localhost:8080")
	router_err := router.Run() // listen and serve on 0.0.0.0:8080
	if router_err != nil {
		log.Fatalf("couldn't run server gg: %v", router_err)
	}
}
