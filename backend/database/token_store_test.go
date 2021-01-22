// Collection of tests for the database access layer of functions evolving
// around tokens.

package database

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/util"
)

var (
	// tToken is a mock token for testing purposes
	tToken = x.Token{
		TokenID: util.GenerateString(43),
		UserID:  1,
		Expiry:  time.Now().Add(time.Hour * 1),
	}
)

func TestGetToken(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TokenStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM tokens"

	table := []string{"token_id", "user_id", "expiry"}

	// Declare test cases
	tests := []struct {
		name      string
		tokenID   string
		mock      func(tokenID string)
		wantToken x.Token
		wantError bool
	}{
		{
			// When everything works as intended
			name:    "#1 OK",
			tokenID: tToken.TokenID,
			mock: func(tokenID string) {
				rows := sqlmock.NewRows(table).
					AddRow(tToken.TokenID, tToken.UserID, tToken.Expiry)

				mock.ExpectQuery(queryMatch).WithArgs(tokenID).WillReturnRows(rows)
			},
			wantToken: tToken,
			wantError: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.tokenID)

			event, err := store.GetToken(test.tokenID)

			if (err != nil) != test.wantError {
				t.Errorf("GetToken() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(event, test.wantToken) {
				t.Errorf("GetEvent() = %v, want %v", event, test.wantToken)
			}
		})
	}
}
