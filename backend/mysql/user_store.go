package mysql

/*
 * user_store.go contains all functions for users that require database access
 */

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

/*
 * UserStore implements database access
 */
type UserStore struct {
	*sqlx.DB
}

/*
 * User gets user by user ID
 */
func (store UserStore) User(userID int) (backend.User, error) {
	var u backend.User
	query := `SELECT * FROM users WHERE user_id = ?`
	if err := store.Get(&u, query, userID); err != nil {
		return backend.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

/*
 * UserByUsername gets user by username
 */
func (store UserStore) UserByUsername(username string) (backend.User, error) {
	var u backend.User
	query := `SELECT * FROM users WHERE username = ?`
	if err := store.Get(&u, query, username); err != nil {
		return backend.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

/*
 * Users gets users sorted by admin and user ID
 */
func (store UserStore) Users() ([]backend.User, error) {
	var uu []backend.User
	query := `SELECT * FROM users ORDER BY admin DESC, user_id` // order by [1.] admin (true -> false), [2.] user_id (101 -> 1)
	if err := store.Select(&uu, query); err != nil {
		return []backend.User{}, fmt.Errorf("error getting topics: %w", err)
	}
	return uu, nil
}

/*
 * UsersCount gets number of users
 */
func (store *UserStore) UsersCount() (int, error) {
	var uCount int
	query := `SELECT COUNT(*) FROM users`
	if err := store.Get(&uCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of users: %w", err)
	}
	return uCount, nil
}

/*
 * CreateUser creates user
 */
func (store UserStore) CreateUser(u *backend.User) error {
	query := `INSERT INTO users(username, password, admin) VALUES (?, ?, ?)`
	if _, err := store.Exec(query, u.Username, u.Password, u.Admin); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

/*
 * UpdateUser updates user
 */
func (store UserStore) UpdateUser(u *backend.User) error {
	query := `UPDATE users SET password = ? WHERE username = ?`
	if _, err := store.Exec(query, u.Password, u.Username); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

/*
 * DeleteUser deletes user by username
 */
func (store UserStore) DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE user_id = ?`
	if _, err := store.Exec(query, userID); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
