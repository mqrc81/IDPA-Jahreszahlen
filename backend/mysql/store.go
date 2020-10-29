package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewStore(dataSourceName string) (*Store, error) {
	// configure database connection
	db, err := sqlx.Open("mysql", dataSourceName) // username:password@host/mysql(address)?param=value
	if err != nil {
		return nil, fmt.Errorf("error opening database: #{err}")
	}

	// initialize database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: #{err}")
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
