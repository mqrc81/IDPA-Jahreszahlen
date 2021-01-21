package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// TokenStore is the database access object.
type TokenStore struct {
	*sqlx.DB
}

// GetToken gets a token by ID.
func (store *TokenStore) GetToken(tokenID string) (x.Token, error) {
	var token x.Token

	// Execute prepared statement
	query := `
		SELECT * 
		FROM tokens 
		WHERE token_id = ?
		`
	if err := store.Get(&token, query, tokenID); err != nil {
		return x.Token{}, fmt.Errorf("error getting token: %w", err)
	}

	return token, nil
}

// CreateToken creates a new token.
func (store *TokenStore) CreateToken(token *x.Token) error {

	// Execute prepared statement
	query := `
		INSERT INTO tokens(token_id, user_id, expiry) 
		VALUES (?, ?, ?)
		`
	if _, err := store.Exec(query,
		token.TokenID,
		token.UserID,
		time.Now().Add(time.Hour)); err != nil {
		return fmt.Errorf("error creating token: %w", err)
	}

	return nil
}

// DeleteTokensByUser deletes all existing tokens of a certain user.
func (store *TokenStore) DeleteTokensByUser(userID int) error {

	// Execute prepared statement
	query := `
		DELETE FROM tokens 
		WHERE user_id = ?
		`
	if _, err := store.Exec(query, userID); err != nil {
		return fmt.Errorf("error deleting tokens: %w", err)
	}

	return nil
}
