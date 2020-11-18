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

	// Establish database connection
	dsn := os.Getenv("DB_DSN")
	store, err := mysql.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	//// Create CSRF protection key
	//csrfKey := []byte("01234567890123456789012345678901")

	// Host website
	handler := web.NewHandler(store/*, csrfKey*/)
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}
