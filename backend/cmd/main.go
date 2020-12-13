package main

/*
 * main.go is the main file, which obtains connection to database and server
 */

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/database"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/web"
)

/*
 * main is the initial starting point of the program
 */
func main() {
	// Access global environment variables
	if err := godotenv.Load("backend/.env"); err != nil {
		log.Fatal(err)
	}

	// Get data-source name from environment variables
	dsn := os.Getenv("DB_DSN")

	// Establish database connection
	store, err := database.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize session manager
	sessions, err := web.NewSessionManager(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Serve website
	handler := web.NewHandler(store, sessions)
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}
