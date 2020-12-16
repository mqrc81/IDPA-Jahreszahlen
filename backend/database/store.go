package database

// store.go
// Pivot of all stores

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// NewStore
// Connects to database and initializes new store objects
func NewStore(dataSourceName string) (*Store, error) {
	// Opens database connection
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Pings database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &Store{
		&TopicStore{DB: db},
		&EventStore{DB: db},
		&UserStore{DB: db},
		&ScoreStore{DB: db},
	}, nil
}

// Store
// Combines all stores
type Store struct {
	*TopicStore
	*EventStore
	*UserStore
	*ScoreStore
}
