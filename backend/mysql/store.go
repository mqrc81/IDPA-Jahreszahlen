package mysql

/*
 * TODO Header
 */

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

/*
 * Connects to MySQL ClearDB database via Heroku.com
 */
func NewStore(dsn string) (*Store, error) {
	// Opens and pings database
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening or pinning database connection: %w", err)
	}

	return &Store {
		&TopicStore{db},
		&EventStore{db},
	}, nil
}

type Store struct {
	*TopicStore
	*EventStore
}
