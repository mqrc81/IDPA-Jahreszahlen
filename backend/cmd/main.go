package main

// main.go
// The main file, which initializes a connection to the database and the server.
// It also obtains session management and CSRF-protection.

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/database"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/web"
)

// main
// The inital starting point of the program, which initializes a connection to
// the database and the server. It also obtains session management and CSRF-
// protection.
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

	// Initializes HTTP-handlers, including router and middleware
	handler := web.NewHandler(store, sessions)

	// Listen on the TCP network address and call Serve with handler to handle
	// requests on incoming connections.
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}
