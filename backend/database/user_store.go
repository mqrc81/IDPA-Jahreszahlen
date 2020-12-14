package database

// user_store.go
// Part of the database layer. Contains all functions for users that access the
// database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// UserStore
// The database access object
type UserStore struct {
	*sqlx.DB
}

// User
// Gets a user by ID
func (store UserStore) User(userID int) (backend.User, error) {
	var u backend.User

	// Execute prepared statement
	query := `SELECT * FROM users WHERE user_id = ?`
	if err := store.Get(&u, query, userID); err != nil {
		return backend.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

// UserByUsername
// Gets a user by username
func (store UserStore) UserByUsername(username string) (backend.User, error) {
	var u backend.User

	// Execute prepared statement
	query := `SELECT * FROM users WHERE username = ?`
	if err := store.Get(&u, query, username); err != nil {
		return backend.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

// Users
// Gets all users
func (store UserStore) Users() ([]backend.User, error) {
	var uu []backend.User

	// Execute prepared statement
	query := `SELECT * FROM users ORDER BY admin DESC, user_id` // order by [1.] admin (true -> false), [2.] user_id (101 -> 1)
	if err := store.Select(&uu, query); err != nil {
		return []backend.User{}, fmt.Errorf("error getting topics: %w", err)
	}
	return uu, nil
}

// CountUsers
// Gets amount of users
func (store *UserStore) CountUsers() (int, error) {
	var uCount int

	// Execute prepared statement
	query := `SELECT COUNT(*) FROM users`
	if err := store.Get(&uCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of users: %w", err)
	}
	return uCount, nil
}

// CreateUser
// Creates a new user
func (store UserStore) CreateUser(u *backend.User) error {
	// Execute prepared statement
	query := `INSERT INTO users(username, password, admin) VALUES (?, ?, ?)`
	if _, err := store.Exec(query, u.Username, u.Password, u.Admin); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

// UpdateUser
// Updates an existing user
func (store UserStore) UpdateUser(u *backend.User) error {
	// Execute prepared statement
	query := `UPDATE users SET password = ? WHERE username = ?`
	if _, err := store.Exec(query, u.Password, u.Username); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

// DeleteUser
// Deletes an existing user
func (store UserStore) DeleteUser(userID int) error {
	// Execute prepared statement
	query := `DELETE FROM users WHERE user_id = ?`
	if _, err := store.Exec(query, userID); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
