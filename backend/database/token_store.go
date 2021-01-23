package database

import (
	"fmt"

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

	query := `
		SELECT * 
		FROM tokens 
		WHERE token_id = ?
		`

	// Execute prepared statement
	if err := store.Get(&token, query, tokenID); err != nil {
		return x.Token{}, fmt.Errorf("error getting token: %w", err)
	}

	return token, nil
}

// CreateToken creates a new token.
func (store *TokenStore) CreateToken(token *x.Token) error {

	query := `
		INSERT INTO tokens(token_id, user_id, expiry) 
		VALUES (?, ?, ?)
		`

	// Execute prepared statement
	if _, err := store.Exec(query,
		token.TokenID,
		token.UserID,
		token.Expiry,
	); err != nil {
		return fmt.Errorf("error creating token: %w", err)
	}

	return nil
}

// DeleteTokensByUser deletes all existing tokens of a certain user.
func (store *TokenStore) DeleteTokensByUser(userID int) error {

	query := `
		DELETE FROM tokens 
		WHERE user_id = ?
		`

	// Execute prepared statement
	if _, err := store.Exec(query, userID); err != nil {
		return fmt.Errorf("error deleting tokens: %w", err)
	}

	return nil
}
