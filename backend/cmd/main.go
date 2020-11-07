package main

import (
	"fmt"
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

	dbURL := fmt.Sprint("mysql://", os.Getenv("DB_USERNAME"), ":", os.Getenv("DB_PASSWORD"), "@", os.Getenv("DB_ADDRESS"), "/", os.Getenv("DB_NAME"), "?reconnect=true")
	fmt.Println(dbURL)
	store, err := mysql.NewStore(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	handler := web.NewHandler(store)
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}
