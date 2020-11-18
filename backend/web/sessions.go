package web

import (
	"database/sql"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

func NewSessionManager(dataSourceName string) (*scs.SessionManager, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	sessions := scs.New()
	sessions.Store = mysqlstore.New(db)
	return sessions, nil
}
