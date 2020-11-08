package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewStore(dsn string) (*Store, error) {
	// configure database connection
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening or pinning database connection: %w", err)
	}

	return &Store{
		&UnitStore{db},
		&EventStore{db},
	}, nil
}

type Store struct {
	*UnitStore
	*EventStore
}
