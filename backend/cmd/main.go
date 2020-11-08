package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/mysql"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/web"
)

func main() {
	if err := godotenv.Load("backend/.env"); err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("DB_DSN")
	store, err := mysql.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	handler := web.NewHandler(store)
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}
