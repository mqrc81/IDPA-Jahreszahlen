// The database store evolving around users, with all necessary methods that
// access the database.

package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// UserStore is the database access object.
type UserStore struct {
	*sqlx.DB
}

// GetUser gets a user by ID.
func (store *UserStore) GetUser(userID int) (x.User, error) {
	var user x.User

	query := `
		SELECT u.*, 
		       COUNT(DISTINCT s.score_id) AS scores_count
		FROM users u 
		    LEFT JOIN scores s ON s.user_id = u.user_id
		WHERE u.user_id = ?
		`

	// Execute prepared statement
	if err := store.Get(&user, query, userID); err != nil {
		return x.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// GetUserByUsername gets a user by its username.
func (store *UserStore) GetUserByUsername(username string) (x.User, error) {
	var user x.User

	query := `
		SELECT u.*, 
		       COUNT(DISTINCT s.score_id) AS scores_count
		FROM users u 
		    LEFT JOIN scores s ON s.user_id = u.user_id
		WHERE u.username = ?
		`

	// Execute prepared statement
	if err := store.Get(&user, query, username); err != nil {
		return x.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// GetUserByEmail gets a user by its email.
func (store *UserStore) GetUserByEmail(email string) (x.User, error) {
	var user x.User

	query := `
		SELECT u.*, 
		       COUNT(DISTINCT s.score_id) AS scores_count
		FROM users u 
		    LEFT JOIN scores s ON s.user_id = u.user_id
		WHERE u.email = ?
		`

	// Execute prepared statement
	if err := store.Get(&user, query, email); err != nil {
		return x.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// GetUsers gets all users.
func (store *UserStore) GetUsers() ([]x.User, error) {
	var users []x.User

	query := `
		SELECT u.*,
		       COUNT(DISTINCT s.score_id) AS scores_count
		FROM users u
		    LEFT JOIN scores s ON s.user_id = u.user_id
		GROUP BY u.user_id, u.admin, u.username
		ORDER BY u.admin DESC, u.username 
		` // Sorted in alphabetical order, but all admins first

	// Execute prepared statement
	if err := store.Select(&users, query); err != nil {
		return []x.User{}, fmt.Errorf("error getting users: %w", err)
	}

	return users, nil
}

// CountUsers gets amount of users.
func (store *UserStore) CountUsers() (int, error) {
	var userCount int

	query := `
		SELECT COUNT(user_id) 
		FROM users
		`

	// Execute prepared statement
	if err := store.Get(&userCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of users: %w", err)
	}

	return userCount, nil
}

// CreateUser creates a new user.
func (store *UserStore) CreateUser(user *x.User) error {

	query := `
		INSERT INTO users(username, email, password, admin) 
		VALUES (?, ?, ?, ?)
		`

	// Execute prepared statement
	if _, err := store.Exec(query,
		user.Username,
		user.Email,
		user.Password,
		user.Admin,
	); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user.
func (store *UserStore) UpdateUser(user *x.User) error {

	query := `
		UPDATE users 
		SET username = ?, email = ?, password = ?, admin = ?, verified = ? 
		WHERE user_id = ?
		`

	// Execute prepared statement
	if _, err := store.Exec(query,
		user.Username,
		user.Email,
		user.Password,
		user.Admin,
		user.Verified,
		user.UserID,
	); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

// DeleteUser deletes an existing user.
func (store *UserStore) DeleteUser(userID int) error {

	query := `
		DELETE FROM users 
		WHERE user_id = ?
		`

	// Execute prepared statement
	if _, err := store.Exec(query, userID); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
