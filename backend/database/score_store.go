package database

// score_store.go
// Part of the database layer. Contains all functions for scores that access the
// database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// ScoreStore
// The database access object
type ScoreStore struct {
	*sqlx.DB
}

// Scores
// Gets a certain amount of scores with a certain offset, sorted by points
// descending
func (store ScoreStore) Scores(limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score

	// Execute prepared statement
	if err := store.Select(&ss, `SELECT * FROM scores ORDER BY points DESC LIMIT ?, ?`,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

// ScoresByTopic
// Gets a certain amount of scores of a certain topic with a certain offset,
// sorted by points descending
func (store ScoreStore) ScoresByTopic(topicID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score

	// Execute prepared statement
	if err := store.Select(&ss, `SELECT * FROM scores WHERE topic_id = ? ORDER BY points DESC LIMIT ?, ?`,
		topicID,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

// ScoresByUser
// Gets a certain amount of scores of a certain user with a certain offset,
// sorted by points descending
func (store ScoreStore) ScoresByUser(userID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score

	// Execute prepared statement
	if err := store.Select(&ss, `SELECT * FROM scores WHERE user_id = ? ORDER BY points DESC LIMIT ?, ?`,
		userID,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

// ScoresByTopicAndUser
// Gets a certain amount of scores of a certain topic and user with a certain
// offset, sorted by points descending
func (store ScoreStore) ScoresByTopicAndUser(topicID int, userID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score

	// Execute prepared statement
	query := `SELECT * FROM scores WHERE topic_id = ? AND user_id = ? ORDER BY points DESC LIMIT ?, ?`
	if err := store.Select(&ss, query, topicID, userID, offset, limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

// CountScores
// Gets amount of scores
func (store *ScoreStore) CountScores() (int, error) {
	var sCount int

	// Execute prepared statement
	query := `SELECT COUNT(*) FROM scores`
	if err := store.Get(&sCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of scores: %w", err)
	}
	return sCount, nil
}

// CreateScore
// Creates a new score
func (store ScoreStore) CreateScore(s *backend.Score) error {
	// Execute prepared statement
	query := `INSERT INTO scores(topic_id, user_id, points, date) VALUES (?, ?, ?, ?)`
	if _, err := store.Exec(query,
		s.TopicID,
		s.UserID,
		s.Points,
		s.Date); err != nil {
		return fmt.Errorf("error creating s: %w", err)
	}

	// Execute prepared statement
	query = `SELECT * FROM scores WHERE score_id = last_insert_id()`
	if err := store.Get(s, query); err != nil {
		return fmt.Errorf("error getting created s: %w", err)
	}
	return nil
}
