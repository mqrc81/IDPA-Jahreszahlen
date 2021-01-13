// The main file, which acquires a connection to the database and the
// server. It also obtains session management and CSRF-protection.

package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/database"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/web"
)

// main is the initial starting point of the program, which acquires a
// connection to the database and the server. It also obtains session
// management and CSRF-protection.
func main() {
	// Access global environment variables
	if err := godotenv.Load("backend/.env"); err != nil {
		log.Fatal(err)
	}

	// Get data-source-name from environment variables
	dataSourceName := os.Getenv("DB_DSN")

	// Establish database connection with the help of the data-source-name
	store, err := database.NewStore(dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize session manager
	sessions, err := web.NewSessionManager(dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Generate random 32-byte key for CSRF-protection
	csrfKey := make([]byte, 32)
	_, err = rand.Read(csrfKey)
	if err != nil {
		log.Fatalf("error generating csrf-protection csrfKey: %e", err)
	}

	// Initialize HTTP-handlers, including router and middleware
	handler := web.NewHandler(store, sessions, csrfKey)

	// Listen on the TCP network address and call Serve with handler to handle
	// requests on incoming connections
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}
