// Collection of tests for the database access layer of functions evolving
// around users.

package database

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	// tUser is a mock user for testing purposes
	tUser = x.User{
		UserID:      1,
		Username:    "user_1",
		Email:       "user1@mail.com",
		Password:    "Passw0rd!",
		Admin:       false,
		Verified:    false,
		ScoresCount: 20,
	}

	// tUser2 is a mock user for testing purposes
	tUser2 = x.User{
		UserID:      2,
		Username:    "user_2",
		Email:       "user2@mail.com",
		Password:    "Passw0rd!",
		Admin:       false,
		Verified:    true,
		ScoresCount: 30,
	}

	// tUser3 is a mock user for testing purposes
	tUser3 = x.User{
		UserID:      3,
		Username:    "admin_1",
		Email:       "admin_1@mail.com",
		Password:    "Passw0rd!",
		Admin:       true,
		Verified:    true,
		ScoresCount: 50,
	}

	// nilUsers is a nil slice of users, since "var u []User" is a nil slice
	// and "u := []User{}" is an empty slice (so we can't use the latter for
	// this use case)
	nilUsers []x.User
)

// TestGetUser tests getting a user by ID.
func TestGetUser(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &UserStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM users"

	table := []string{"user_id", "username", "email", "password", "admin", "verified", "scores_count"}

	// Declare test cases
	tests := []struct {
		name      string
		userID    int
		mock      func(userID int)
		wantUser  x.User
		wantError bool
	}{
		{
			// When everything works as intended
			name:   "#1 OK",
			userID: tUser.UserID,
			mock: func(userID int) {
				rows := sqlmock.NewRows(table).
					AddRow(tUser.UserID, tUser.Username, tUser.Email, tUser.Password, tUser.Admin, tUser.Verified,
						tUser.ScoresCount)

				mock.ExpectQuery(queryMatch).WithArgs(userID).WillReturnRows(rows)
			},
			wantUser:  tUser,
			wantError: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.userID)

			user, err := store.GetUser(test.userID)

			if (err != nil) != test.wantError {
				t.Errorf("GetUser() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("GetUser() = %v, want %v", user, test.wantUser)
			}

		})
	}
}

// TestGetUserByUsername tests getting a user by its username.
func TestGetUserByUsername(t *testing.T) {

}

// TestGetUserByEmail tests getting a user by its email.
func TestGetUserByEmail(t *testing.T) {

}

// TestGetUsers tests getting all users.
func TestGetUsers(t *testing.T) {

}

// TestCountUsers tests getting amount of users.
func TestCountUsers(t *testing.T) {

}

// TestCreateUser tests creating a new user.
func TestCreateUser(t *testing.T) {

}

// TestUpdateUser tests updating an existing user.
func TestUpdateUser(t *testing.T) {

}

// TestDeleteUser tests deleting an existing user.
func TestDeleteUser(t *testing.T) {

}
