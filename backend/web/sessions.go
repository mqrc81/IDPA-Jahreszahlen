package web

// sessions.go
// Contains session management.

import (
	"context"
	"database/sql"
	"encoding/gob"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// init
// Gets initialized with the package. Registers certain types to the session,
// because by default the session can only contain basic data types (int, bool,
// string, etc.).
func init() {
	gob.Register(backend.Event{})
	gob.Register([]backend.Event{})
}

// NewSessionManager
// Initializes new session management.
func NewSessionManager(dataSourceName string) (*scs.SessionManager, error) {
	// Opens MySQL connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Creates new sessions
	sessions := scs.New()
	sessions.Store = mysqlstore.New(db)
	return sessions, nil
}

// SessionData
// Holds data to be accessed through the session.
type SessionData struct {
	FlashMessageSuccess string
	FlashMessageError   string
	Form                interface{}
	User                backend.User
	LoggedIn            bool
	PlayEvents          []backend.Event
}

// GetSessionData
// Gets all the data from the session
func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	// Retrieve flash message from session
	data.FlashMessageSuccess = session.PopString(ctx, "flash_success")
	data.FlashMessageError = session.PopString(ctx, "flash_error")

	// Retrieve form from session
	data.Form = session.Pop(ctx, "form")
	if data.Form == nil {
		data.Form = map[string]string{}
	}

	// Retrieve user from session
	userInf := ctx.Value("user")
	if userInf != nil { // 'If there is a user in the session'
		data.User = userInf.(backend.User)
		data.LoggedIn = true
	} else {
		data.User = backend.User{}
		data.LoggedIn = false
	}

	// Retrieve events from session
	data.PlayEvents = session.Pop(ctx, "events").([]backend.Event)
	if data.PlayEvents == nil {
		data.PlayEvents = []backend.Event{}
	}

	return data
}
