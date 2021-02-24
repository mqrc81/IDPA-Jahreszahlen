// The main file, which acquires a connection to the database and the
// server. It also obtains session management and CSRF-protection.

package main

import (
	"crypto/rand"
	"fmt"
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
	fmt.Println("Starting application...")

	// Access global environment variables
	// If this is a local environment, the variable will be empty, since it
	// didn't load the .env file yet
	// If this is production environment (hosted on heroku), then it shouldn't
	// load any file, which would result in crashing the application because of
	// the error being produced
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("error loading environment variables: %v", err)
		}
	}

	// Get data-source-name from environment variables
	dataSourceName := os.Getenv("MYSQL_DSN")

	// Establish database connection with the help of the data-source-name
	store, err := database.NewStore(dataSourceName)
	if err != nil {
		log.Fatalf("error initializing new database store: %v", err)
	}

	// Initialize session manager
	sessions, err := web.NewSessionManager(dataSourceName)
	if err != nil {
		log.Fatalf("error initializing new session manager: %v", err)
	}

	// Generate random 32-byte key for CSRF-protection
	csrfKey := make([]byte, 32)
	if _, err = rand.Read(csrfKey); err != nil {
		log.Fatalf("error generating csrf-key: %v", err)
	}

	// Initialize HTTP-handlers, including router and middleware
	handler := web.NewHandler(store, sessions, csrfKey)

	// Listen on the TCP network address and call Serve with handler to handle
	// requests on incoming connections
	port := ":" + os.Getenv("PORT")
	fmt.Println("Listening on port " + port + "...")
	if err = http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("error listening on the tcp network: %v", err)
	}
}
