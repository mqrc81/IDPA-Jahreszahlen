// The pivot of all database stores, which is responsible for initializing a
// database connection and combining all existing store objects into one single
// store object to be accesses throughout the HTTP-handlers.

package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// NewStore connects to database and initializes new store objects
func NewStore(dataSourceName string) (*Store, error) {
	// Open database connection
	db, err := sqlx.Open("mysql", dataSourceName+"?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Ping database connection
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

// Store combines all stores.
type Store struct {
	*TopicStore
	*EventStore
	*UserStore
	*ScoreStore
}
