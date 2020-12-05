package web

/*
 * sessions.go adds sessions-management
 */

import (
	"context"
	"database/sql"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
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
	Form         interface{}
	User         backend.User
	LoggedIn     bool
	Admin        bool
}

/*
 * GetSessionData gets session data
 */
func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	// Retrieve flash message from session
	data.FlashMessage = session.PopString(ctx, "flash")

	// Retrieve form from session
	data.Form = session.PopInt(ctx, "form")
	if data.Form == nil {
		data.Form = map[string]string{}
	}

	// Retrieve user from session
	userInf := ctx.Value("user")
	if userInf != nil {
		data.User = userInf.(backend.User)
		data.LoggedIn = false
	} else {
		data.User = backend.User{}
		data.LoggedIn = true
	}

	data.Admin = data.User.Admin

	return data
}
