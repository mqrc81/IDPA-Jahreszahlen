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
	var user backend.User

	// Execute prepared statement
	query := `SELECT * FROM users WHERE user_id = ?`
	if err := store.Get(&user, query, userID); err != nil {
		return backend.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

// UserByUsername
// Gets a user by username
func (store UserStore) UserByUsername(username string) (backend.User, error) {
	var user backend.User

	// Execute prepared statement
	query := `SELECT * FROM users WHERE username = ?`
	if err := store.Get(&user, query, username); err != nil {
		return backend.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

// Users
// Gets all users
func (store UserStore) Users() ([]backend.User, error) {
	var users []backend.User

	// Execute prepared statement
	query := `SELECT * FROM users ORDER BY admin DESC, username` // Sorted in alphabetical order, admins first
	if err := store.Select(&users, query); err != nil {
		return []backend.User{}, fmt.Errorf("error getting topics: %w", err)
	}
	return users, nil
}

// CountUsers
// Gets amount of users
func (store *UserStore) CountUsers() (int, error) {
	var userCount int

	// Execute prepared statement
	query := `SELECT COUNT(*) FROM users`
	if err := store.Get(&userCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of users: %w", err)
	}
	return userCount, nil
}

// CreateUser
// Creates a new user
func (store UserStore) CreateUser(user *backend.User) error {
	// Execute prepared statement
	query := `INSERT INTO users(username, password, admin) VALUES (?, ?, ?)`
	if _, err := store.Exec(query,
		user.Username,
		user.Password,
		user.Admin); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

// UpdateUser
// Updates an existing user
func (store UserStore) UpdateUser(user *backend.User) error {
	// Execute prepared statement
	query := `UPDATE users SET password = ?, username = ?, admin = ? WHERE user_id = ?`
	if _, err := store.Exec(query,
		user.Password,
		user.Username,
		user.Admin,
		user.UserID); err != nil {
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
