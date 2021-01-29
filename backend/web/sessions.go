// Contains session management, which is responsible for authentication and
// verification of users and transporting forms, flash messages and other
// structures such as quiz data.

package web

import (
	"context"
	"database/sql"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// NewSessionManager initializes new session management.
func NewSessionManager(dataSourceName string) (*scs.SessionManager, error) {

	// Open MySQL connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Configure database connections
	db.SetConnMaxLifetime(time.Minute * 4)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(25)

	// Create new sessions
	sessions := scs.New()
	sessions.Store = mysqlstore.New(db)
	return sessions, nil
}

// SessionData holds data to be accessed through the session.
type SessionData struct {
	FlashMessageSuccess string
	FlashMessageError   string
	Form                interface{}
	User                x.User
	LoggedIn            bool
}

// GetSessionData gets all the data from session.
func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	// Retrieve flash message from session
	data.FlashMessageSuccess = session.PopString(ctx, "flash_success")
	data.FlashMessageError = session.PopString(ctx, "flash_error")

	// Retrieve form from session
	data.Form = session.Pop(ctx, "form")

	// Retrieve user from session
	userInf := ctx.Value("user")
	if userInf != nil { // 'If there is a user in the session'
		data.User = userInf.(x.User)
		data.LoggedIn = true
	} else {
		data.User = x.User{}
		data.LoggedIn = false
	}

	return data
}
