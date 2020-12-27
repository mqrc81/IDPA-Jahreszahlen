package database

/*
 * Part of the database layer. Contains all functions for scores that access
 * the database.
 */

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// ScoreStore is the database access object.
type ScoreStore struct {
	*sqlx.DB
}

// GetScores gets a certain amount of scores with a certain offset, sorted by
// points descending.
func (store ScoreStore) GetScores(limit int, offset int) ([]backend.Score, error) {
	var scores []backend.Score

	// Execute prepared statement
	query := `
		SELECT s.*, 
		       t.name AS topic_name, 
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

// GetScoresByTopic gets a certain amount of scores of a certain topic with a
// certain offset, sorted by points descending.
func (store ScoreStore) GetScoresByTopic(topicID int, limit int, offset int) ([]backend.Score, error) {
	var scores []backend.Score

	// Execute prepared statement
	query := `
		SELECT s.score_id, s.topic_id, s.user_id, s.points, s.date, 
		       t.name AS topic_name, 
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

// GetScoresByUser gets a certain amount of scores of a certain user with a
// certain offset, sorted by points descending.
func (store ScoreStore) GetScoresByUser(userID int, limit int, offset int) ([]backend.Score, error) {
	var scores []backend.Score

	// Execute prepared statement
	query := `
		SELECT s.*, 
		       t.name AS topic_name, 
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

// GetScoresByTopicAndUser gets a certain amount of scores of a certain topic
// and user with a certain offset, sorted by points descending.
func (store ScoreStore) GetScoresByTopicAndUser(topicID int, userID int, limit int, offset int) ([]backend.Score, error) {
	var scores []backend.Score

	// Execute prepared statement
	query := `
		SELECT s.*, 
		       t.name AS topic_name, 
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

// CreateScore creates a new score.
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
