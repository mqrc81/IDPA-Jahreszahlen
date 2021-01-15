// The database store evolving around users, with all necessary methods that
// access the database.

package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/jahreszahlen"
)

// UserStore is the database access object
type UserStore struct {
	*sqlx.DB
}

// GetUser gets a user by ID.
func (store UserStore) GetUser(userID int) (jahreszahlen.User, error) {
	var user jahreszahlen.User

	// Execute prepared statement
	query := `
		SELECT u.*, 
		       COUNT(DISTINCT s.score_id) AS scores_count
		FROM users u 
		    LEFT JOIN scores s ON s.user_id = u.user_id
		WHERE u.user_id = ?
		`
	if err := store.Get(&user, query, userID); err != nil {
		return jahreszahlen.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// GetUserByUsername gets a user by its username.
func (store UserStore) GetUserByUsername(username string) (jahreszahlen.User, error) {
	var user jahreszahlen.User

	// Execute prepared statement
	query := `
		SELECT u.*, 
		       COUNT(DISTINCT s.score_id) AS scores_count
		FROM users u 
		    LEFT JOIN scores s ON s.user_id = u.user_id
		WHERE u.username = ?
		`
	if err := store.Get(&user, query, username); err != nil {
		return jahreszahlen.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// GetUserByEmail gets a user by its email.
func (store UserStore) GetUserByEmail(email string) (jahreszahlen.User, error) {
	var user jahreszahlen.User

	// Execute prepared statement
	query := `
		SELECT u.*, 
		       COUNT(DISTINCT s.score_id) AS scores_count
		FROM users u 
		    LEFT JOIN scores s ON s.user_id = u.user_id
		WHERE u.email = ?
		`
	if err := store.Get(&user, query, email); err != nil {
		return jahreszahlen.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// GetUsers gets all users.
func (store UserStore) GetUsers() ([]jahreszahlen.User, error) {
	var users []jahreszahlen.User

	// Execute prepared statement
	query := `
		SELECT u.*,
		       COUNT(DISTINCT s.score_id) AS scores_count
		FROM users u
		    LEFT JOIN scores s ON s.user_id = u.user_id
		GROUP BY u.user_id, u.admin, u.username
		ORDER BY u.admin DESC, u.username   
		` // Sorted in alphabetical order, but all admins first
	if err := store.Select(&users, query); err != nil {
		return []jahreszahlen.User{}, fmt.Errorf("error getting topics: %w", err)
	}

	return users, nil
}

// CountUsers gets amount of users.
func (store *UserStore) CountUsers() (int, error) {
	var userCount int

	// Execute prepared statement
	query := `
		SELECT COUNT(*) 
		FROM users
		`
	if err := store.Get(&userCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of users: %w", err)
	}

	return userCount, nil
}

// CreateUser creates a new user.
func (store UserStore) CreateUser(user *jahreszahlen.User) error {
	// Execute prepared statement
	query := `
		INSERT INTO users(username, password, admin) 
		VALUES (?, ?, ?)
		`
	if _, err := store.Exec(query,
		user.Username,
		user.Password,
		user.Admin); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user.
func (store UserStore) UpdateUser(user *jahreszahlen.User) error {
	// Execute prepared statement
	query := `
		UPDATE users 
		SET password = ?, username = ?, admin = ? 
		WHERE user_id = ?
		`
	if _, err := store.Exec(query,
		user.Password,
		user.Username,
		user.Admin,
		user.UserID); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

// DeleteUser deletes an existing user.
func (store UserStore) DeleteUser(userID int) error {
	// Execute prepared statement
	query := `
		DELETE FROM users 
		WHERE user_id = ?
		`
	if _, err := store.Exec(query, userID); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
