package main

/*
 * TODO Header
 */

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/mysql"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/web"
)

/*
 * Main function
 */
func main() {
	// Load '.env' file, where global environment variables are stored
	if err := godotenv.Load("backend/.env"); err != nil {
		log.Fatal(err)
	}
	dsn := os.Getenv("DB_DSN")

	// Establish database connection
	store, err := mysql.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize session manager
	sessions, err := web.NewSessionManager(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Host website
	handler := web.NewHandler(store, sessions)
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}
