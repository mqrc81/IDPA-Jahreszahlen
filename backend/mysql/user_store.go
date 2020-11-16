package mysql

/*
 * TODO Header
 */

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type UserStore struct {
	*sqlx.DB
}

/*
 * User gets user by username
 */
func (s UserStore) User(username string) (backend.User, error) {
	var u backend.User
	if err := s.Get(&u, `SELECT * FROM users WHERE username = ?`, username); err != nil {
		return backend.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

/*
 * CreateUser creates user
 */
func (s UserStore) CreateUser(user *backend.User) error {
	if _, err := s.Exec(`INSERT INTO users(username, password, admin) VALUES (?, ?, ?)`,
		user.Username,
		user.Password,
		user.Admin); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

/*
 * UpdateUser updates user
 */
func (s UserStore) UpdateUser(user *backend.User) error {
	if _, err := s.Exec(`UPDATE users SET password = ? WHERE username = ?`,
		user.Password,
		user.Username); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

/*
 * DeleteUser deletes user by username
 */
func (s UserStore) DeleteUser(username string) error {
	if _, err := s.Exec(`DELETE FROM users WHERE username = ?`, username); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
