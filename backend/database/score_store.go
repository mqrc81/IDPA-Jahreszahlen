package database

/*
 * Part of the database layer. Contains all functions for scores that access
 * the database.
 */

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// ScoreStore is the database access object.
type ScoreStore struct {
	*sqlx.DB
}

// GetScores gets all scores, sorted by points descending.
func (store ScoreStore) GetScores() ([]backend.Score, error) {
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
		`
	if err := store.Select(&scores, query); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}

	return scores, nil
}

// GetScoresByTopic gets scores of a certain topic, sorted by points
// descending.
func (store ScoreStore) GetScoresByTopic(topicID int) ([]backend.Score, error) {
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
		`
	if err := store.Select(&scores, query, topicID); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}

	return scores, nil
}

// GetScoresByUser gets scores of a certain user, sorted by points descending.
func (store ScoreStore) GetScoresByUser(userID int) ([]backend.Score, error) {
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
		`
	if err := store.Select(&scores, query, userID); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}

	return scores, nil
}

// GetScoresByTopicAndUser gets scores of a certain topic and user, sorted by
// points descending.
func (store ScoreStore) GetScoresByTopicAndUser(topicID int, userID int) ([]backend.Score, error) {
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
		`
	if err := store.Select(&scores, query, topicID, userID); err != nil {
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
		time.Now().Format("2006-01-02")); err != nil { // current date formatted as 'yyyy-mm-dd'
		return fmt.Errorf("error creating score: %w", err)
	}

	return nil
}
