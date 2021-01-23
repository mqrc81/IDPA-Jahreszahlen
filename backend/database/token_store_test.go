// Collection of tests for the database access layer of functions evolving
// around tokens.

package database

import (
	"errors"
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
		Expiry:  time.Now().Add(time.Hour),
	}
)

// TestGetToken tests getting a token by ID.
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
		{
			// When token with given token ID doesn't exist
			name:    "#2 NOT FOUND",
			tokenID: "123",
			mock: func(tokenID string) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(tokenID).WillReturnRows(rows)
			},
			wantToken: x.Token{},
			wantError: true,
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

// TestCreateToken tests creating a new token.
func TestCreateToken(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TokenStore{DB: db}
	defer db.Close()

	queryMatch := "INSERT INTO tokens"

	// Declare test cases
	tests := []struct {
		name      string
		token     x.Token
		mock      func(token x.Token)
		wantError bool
	}{
		{
			// When everything works as intended
			name:  "#1 OK",
			token: tToken,
			mock: func(token x.Token) {
				mock.ExpectExec(queryMatch).WithArgs(token.TokenID, token.UserID, token.Expiry).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			// When token with given token ID doesn't exist
			name: "#2 NOT FOUND",
			token: x.Token{
				TokenID: "123",
				UserID:  tToken.UserID,
				Expiry:  tToken.Expiry,
			},
			mock: func(token x.Token) {
				mock.ExpectExec(queryMatch).WithArgs(token.TokenID, token.UserID, token.Expiry).
					WillReturnError(errors.New("token with given id does not exist"))
			},
			wantError: true,
		},
		{
			// When user with given user ID doesn't exist
			name: "#3 USER NOT FOUND",
			token: x.Token{
				TokenID: tToken.TokenID,
				UserID:  0,
				Expiry:  tToken.Expiry,
			},
			mock: func(token x.Token) {
				mock.ExpectExec(queryMatch).WithArgs(token.TokenID, token.UserID, token.Expiry).
					WillReturnError(errors.New("user with given id does not exist"))
			},
			wantError: true,
		},
		{
			// When expiry is missing
			name: "#4 EXPIRY MISSING",
			token: x.Token{
				TokenID: tToken.TokenID,
				UserID:  tToken.UserID,
			},
			mock: func(token x.Token) {
				mock.ExpectExec(queryMatch).WithArgs(token.TokenID, token.UserID, token.Expiry).
					WillReturnError(errors.New("expiry can not be empty"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.token)

			err := store.CreateToken(&test.token)

			if (err != nil) != test.wantError {
				t.Errorf("CreateToken() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}

// TestDeleteTokensByUser tests deleting all tokens of a certain user.
func TestDeleteTokensByUser(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TokenStore{DB: db}
	defer db.Close()

	queryMatch := "DELETE FROM tokens"

	// Declare test cases
	tests := []struct {
		name      string
		userID    int
		mock      func(userID int)
		wantError bool
	}{
		{
			// When everything works as intended
			name:   "#1 OK",
			userID: tToken.UserID,
			mock: func(userID int) {
				mock.ExpectExec(queryMatch).WithArgs(userID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			// When token with given token ID doesn't exist
			name:   "#2 USER NOT FOUND",
			userID: tToken.UserID,
			mock: func(userID int) {
				mock.ExpectExec(queryMatch).WithArgs(userID).
					WillReturnError(errors.New("user with given id does not exist"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.userID)

			err := store.DeleteTokensByUser(test.userID)

			if (err != nil) != test.wantError {
				t.Errorf("DeleteTokensByUser() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}
