// Collection of tests for the database access layer of functions evolving
// around users.

package database

import (
	"errors"
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

	// tUsers is a mock array of users for testing purposes
	tUsers = []x.User{
		tUser,
		{
			UserID:      2,
			Username:    "user_2",
			Email:       "user2@mail.com",
			Password:    "Passw0rd!",
			Admin:       false,
			Verified:    true,
			ScoresCount: 30,
		},
		{
			UserID:      3,
			Username:    "admin_1",
			Email:       "admin_1@mail.com",
			Password:    "Passw0rd!",
			Admin:       true,
			Verified:    true,
			ScoresCount: 50,
		},
	}

	// nilUsers is a nil slice of users, since "var u []User" is a nil slice
	// but "u := []User{}" is an empty slice (so we can't use the latter for
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
		{
			// When user with given user ID doesn't exist
			name:   "#2 NOT FOUND",
			userID: 0,
			mock: func(userID int) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(userID).WillReturnRows(rows)
			},
			wantUser:  x.User{},
			wantError: true,
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

	// New mock database
	db, mock := NewMock()
	store := &UserStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM users"

	table := []string{"user_id", "username", "email", "password", "admin", "verified", "scores_count"}

	// Declare test cases
	tests := []struct {
		name      string
		username  string
		mock      func(username string)
		wantUser  x.User
		wantError bool
	}{
		{
			// When everything works as intended
			name:     "#1 OK",
			username: tUser.Username,
			mock: func(username string) {
				rows := sqlmock.NewRows(table).
					AddRow(tUser.UserID, tUser.Username, tUser.Email, tUser.Password, tUser.Admin, tUser.Verified,
						tUser.ScoresCount)

				mock.ExpectQuery(queryMatch).WithArgs(username).WillReturnRows(rows)
			},
			wantUser:  tUser,
			wantError: false,
		},
		{
			// When user with given username doesn't exist
			name:     "#2 USERNAME NOT FOUND",
			username: "",
			mock: func(username string) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(username).WillReturnRows(rows)
			},
			wantUser:  x.User{},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.username)

			user, err := store.GetUserByUsername(test.username)

			if (err != nil) != test.wantError {
				t.Errorf("GetUserByUsername() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("GetUserByUsername() = %v, want %v", user, test.wantUser)
			}
		})
	}
}

// TestGetUserByEmail tests getting a user by its email.
func TestGetUserByEmail(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &UserStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM users"

	table := []string{"user_id", "username", "email", "password", "admin", "verified", "scores_count"}

	// Declare test cases
	tests := []struct {
		name      string
		email     string
		mock      func(email string)
		wantUser  x.User
		wantError bool
	}{
		{
			// When everything works as intended
			name:  "#1 OK",
			email: tUser.Email,
			mock: func(email string) {
				rows := sqlmock.NewRows(table).
					AddRow(tUser.UserID, tUser.Username, tUser.Email, tUser.Password, tUser.Admin, tUser.Verified,
						tUser.ScoresCount)

				mock.ExpectQuery(queryMatch).WithArgs(email).WillReturnRows(rows)
			},
			wantUser:  tUser,
			wantError: false,
		},
		{
			// When user with given email doesn't exist
			name:  "#2 EMAIL NOT FOUND",
			email: "",
			mock: func(email string) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(email).WillReturnRows(rows)
			},
			wantUser:  x.User{},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.email)

			user, err := store.GetUserByEmail(test.email)

			if (err != nil) != test.wantError {
				t.Errorf("GetUserByEmail() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("GetUserByEmail() = %v, want %v", user, test.wantUser)
			}
		})
	}
}

// TestGetUsers tests getting all users.
func TestGetUsers(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &UserStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM users"

	table := []string{"user_id", "username", "email", "password", "admin", "verified", "scores_count"}

	// Declare test cases
	tests := []struct {
		name      string
		mock      func()
		wantUsers []x.User
		wantError bool
	}{
		{
			// When everything works as intended
			name: "#1 OK",
			mock: func() {
				rows := sqlmock.NewRows(table)
				for _, user := range tUsers {
					rows = rows.AddRow(user.UserID, user.Username, user.Email, user.Password, user.Admin,
						user.Verified, user.ScoresCount)
				}

				mock.ExpectQuery(queryMatch).WillReturnRows(rows)
			},
			wantUsers: tUsers,
			wantError: false,
		},
		{
			// When users table is empty
			name: "#2 OK (NO ROWS)",
			mock: func() {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WillReturnRows(rows)
			},
			wantUsers: nilUsers,
			wantError: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock()

			scores, err := store.GetUsers()

			if (err != nil) != test.wantError {
				t.Errorf("GetUsers() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(scores, test.wantUsers) {
				t.Errorf("GetUsers() = %v, want %v", scores, test.wantError)
			}
		})
	}
}

// TestCountUsers tests getting amount of users.
func TestCountUsers(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &UserStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT COUNT((.+)) FROM users"

	table := []string{"COUNT(*)"}

	// Declare test cases
	tests := []struct {
		name           string
		mock           func()
		wantUsersCount int
		wantError      bool
	}{
		{
			// When everything works as intended
			name: "#1 OK",
			mock: func() {
				rows := sqlmock.NewRows(table).AddRow(3)

				mock.ExpectQuery(queryMatch).WillReturnRows(rows)
			},
			wantUsersCount: 3,
			wantError:      false,
		},
		{
			// When the users table is empty
			name: "#2 NO ROWS",
			mock: func() {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WillReturnRows(rows).WillReturnError(errors.New("no users found"))
			},
			wantUsersCount: 0,
			wantError:      true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock()

			usersCount, err := store.CountUsers()

			if (err != nil) != test.wantError {
				t.Errorf("CountUsers() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(usersCount, test.wantUsersCount) {
				t.Errorf("CountUsers() = %v, want %v", usersCount, test.wantUsersCount)
			}
		})
	}
}

// TestCreateUser tests creating a new user.
func TestCreateUser(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &UserStore{DB: db}
	defer db.Close()

	queryMatch := "INSERT INTO users"

	// Declare test cases
	tests := []struct {
		name      string
		user      x.User
		mock      func(user x.User)
		wantError bool
	}{
		{
			// When everything works as intended
			name: "#1 OK",
			user: tUser,
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin).
					WillReturnResult(sqlmock.NewResult(int64(user.UserID), 1))
			},
			wantError: false,
		},
		{
			// When username is missing
			name: "#2 USERNAME MISSING",
			user: x.User{
				Email:    tUser.Email,
				Password: tUser.Password,
				Admin:    tUser.Admin,
			},
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin).
					WillReturnError(errors.New("username can not be empty"))
			},
			wantError: true,
		},
		{
			// When email is missing
			name: "#3 EMAIL MISSING",
			user: x.User{
				Username: tUser.Username,
				Password: tUser.Password,
				Admin:    tUser.Admin,
			},
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin).
					WillReturnError(errors.New("email can not be empty"))
			},
			wantError: true,
		},
		{
			// When password is missing
			name: "#4 PASSWORD MISSING",
			user: x.User{
				Username: tUser.Username,
				Email:    tUser.Email,
				Admin:    tUser.Admin,
			},
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin).
					WillReturnError(errors.New("password can not be empty"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.user)

			err := store.CreateUser(&test.user)

			if (err != nil) != test.wantError {
				t.Errorf("CreateUser() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}

// TestUpdateUser tests updating an existing user.
func TestUpdateUser(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &UserStore{DB: db}
	defer db.Close()

	queryMatch := "UPDATE users"

	// Declare test cases
	tests := []struct {
		name      string
		user      x.User
		mock      func(user x.User)
		wantError bool
	}{
		{
			// When everything works as intended
			name: "#1 OK",
			user: tUser,
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin,
					user.Verified, user.UserID).
					WillReturnResult(sqlmock.NewResult(int64(user.UserID), 1))
			},
			wantError: false,
		},
		{
			// When user with given user ID doesn't exist
			name: "#2 NOT FOUND",
			user: x.User{
				UserID:   0,
				Username: tUser.Username,
				Email:    tUser.Email,
				Password: tUser.Password,
				Admin:    tUser.Admin,
				Verified: tUser.Verified,
			},
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin,
					user.Verified, user.UserID).WillReturnError(errors.New("user with given id does not exist"))
			},
			wantError: true,
		},
		{
			// When username is missing
			name: "#3 USERNAME MISSING",
			user: x.User{
				UserID:   tUser.UserID,
				Email:    tUser.Email,
				Password: tUser.Password,
				Admin:    tUser.Admin,
				Verified: tUser.Verified,
			},
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin,
					user.Verified, user.UserID).
					WillReturnError(errors.New("username can not be empty"))
			},
			wantError: true,
		},
		{
			// When email is missing
			name: "#4 EMAIL MISSING",
			user: x.User{
				UserID:   tUser.UserID,
				Username: tUser.Username,
				Password: tUser.Password,
				Admin:    tUser.Admin,
				Verified: tUser.Verified,
			},
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin,
					user.Verified, user.UserID).
					WillReturnError(errors.New("email can not be empty"))
			},
			wantError: true,
		},
		{
			// When password is missing
			name: "#5 PASSWORD MISSING",
			user: x.User{
				UserID:   tUser.UserID,
				Username: tUser.Username,
				Email:    tUser.Email,
				Admin:    tUser.Admin,
				Verified: tUser.Verified,
			},
			mock: func(user x.User) {
				mock.ExpectExec(queryMatch).WithArgs(user.Username, user.Email, user.Password, user.Admin,
					user.Verified, user.UserID).
					WillReturnError(errors.New("password can not be empty"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.user)

			err := store.UpdateUser(&test.user)

			if (err != nil) != test.wantError {
				t.Errorf("UpdateUser() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}

// TestDeleteUser tests deleting an existing user.
func TestDeleteUser(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &UserStore{DB: db}
	defer db.Close()

	queryMatch := "DELETE FROM users"

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
			userID: tUser.UserID,
			mock: func(userID int) {
				mock.ExpectExec(queryMatch).WithArgs(userID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			// When user with given user ID doesn't exist
			name:   "#2 NOT FOUND",
			userID: 0,
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

			err := store.DeleteUser(test.userID)

			if (err != nil) != test.wantError {
				t.Errorf("DeleteUser() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}
