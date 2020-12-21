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
	var scores []backend.Score

	// Execute prepared statement
	query := `
		SELECT s.score_id, s.topic_id, s.user_id, s.points, s.date, 
		       t.title AS topic_name, 
		       u.username AS user_name
		FROM scores s 
		    LEFT JOIN topics t ON t.topic_id = s.topic_id 
		    LEFT JOIN users u ON u.user_id = s.user_id
		ORDER BY points DESC 
		LIMIT ?, ?
		`
	if err := store.Select(&scores, query,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return scores, nil
}

// ScoresByTopic
// Gets a certain amount of scores of a certain topic with a certain offset,
// sorted by points descending
func (store ScoreStore) ScoresByTopic(topicID int, limit int, offset int) ([]backend.Score, error) {
	var scores []backend.Score

	// Execute prepared statement
	query := `
		SELECT s.score_id, s.topic_id, s.user_id, s.points, s.date, 
		       t.title AS topic_name, 
		       u.username AS user_name
		FROM scores s 
		    LEFT JOIN topics t ON t.topic_id = s.topic_id 
		    LEFT JOIN users u ON u.user_id = s.user_id
		WHERE s.topic_id = ?
		ORDER BY points DESC 
		LIMIT ?, ?
		`
	if err := store.Select(&scores, query,
		topicID,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return scores, nil
}

// ScoresByUser
// Gets a certain amount of scores of a certain user with a certain offset,
// sorted by points descending
func (store ScoreStore) ScoresByUser(userID int, limit int, offset int) ([]backend.Score, error) {
	var scores []backend.Score

	// Execute prepared statement
	query := `
		SELECT s.score_id, s.topic_id, s.user_id, s.points, s.date, 
		       t.title AS topic_name, 
		       u.username AS user_name
		FROM scores s 
		    LEFT JOIN topics t ON t.topic_id = s.topic_id 
		    LEFT JOIN users u ON u.user_id = s.user_id
		WHERE s.user_id = ?
		ORDER BY points DESC 
		LIMIT ?, ?
		`
	if err := store.Select(&scores, query,
		userID,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return scores, nil
}

// ScoresByTopicAndUser
// Gets a certain amount of scores of a certain topic and user with a certain
// offset, sorted by points descending
func (store ScoreStore) ScoresByTopicAndUser(topicID int, userID int, limit int, offset int) ([]backend.Score, error) {
	var scores []backend.Score

	// Execute prepared statement
	query := `
		SELECT s.score_id, s.topic_id, s.user_id, s.points, s.date, 
		       t.title AS topic_name, 
		       u.username AS user_name
		FROM scores s 
		    LEFT JOIN topics t ON t.topic_id = s.topic_id 
		    LEFT JOIN users u ON u.user_id = s.user_id
		WHERE s.topic_id = ? 
		  AND s.user_id = ?
		ORDER BY points DESC 
		LIMIT ?, ?
		`
	if err := store.Select(&scores, query,
		topicID,
		userID,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return scores, nil
}

// CreateScore
// Creates a new score
func (store ScoreStore) CreateScore(score *backend.Score) error {
	// Execute prepared statement
	query := `
		INSERT INTO scores(topic_id, user_id, points, date) 
		VALUES (?, ?, ?, ?)
		`
	if _, err := store.Exec(query,
		score.TopicID,
		score.UserID,
		score.Points,
		score.Date); err != nil {
		return fmt.Errorf("error creating score: %w", err)
	}

	return nil
}
