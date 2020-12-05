package web

/*
 * sessions.go adds sessions-management
 */

import (
	"context"
	"database/sql"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

/*
 * NewSessionManager creates new session
 */
func NewSessionManager(dataSourceName string) (*scs.SessionManager, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	sessions := scs.New()
	sessions.Store = mysqlstore.New(db)
	return sessions, nil
}

/*
 * SessionData hold session data
 */
type SessionData struct {
	FlashMessage string
	User         int
	Form         interface{}
}

/*
 * GetSessionData gets session data
 */
func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	//data.FlashMessage = session.PopString(ctx, "flash")
	//data.User = session.PopInt(ctx, "user")
	data.Form = session.PopInt(ctx, "form")
	if data.Form == nil {
		data.Form = map[string]string{}
	}

	return data
}
