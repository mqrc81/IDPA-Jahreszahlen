// The database store evolving around scores, with all necessary methods that
// access the database.

package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// ScoreStore is the database access object.
type ScoreStore struct {
	*sqlx.DB
}

// GetScores gets all scores, sorted by points descending.
func (store *ScoreStore) GetScores() ([]x.Score, error) {
	var scores []x.Score

	query := `
		SELECT s.*, 
		       t.name AS topic_name, 
		       u.username AS user_name
		FROM scores s 
		    LEFT JOIN topics t ON t.topic_id = s.topic_id 
		    LEFT JOIN users u ON u.user_id = s.user_id
		ORDER BY points DESC
		`

	// Execute prepared statement
	if err := store.Select(&scores, query); err != nil {
		return []x.Score{}, fmt.Errorf("error getting scores: %w", err)
	}

	return scores, nil
}

// GetScoresByTopic gets scores of a certain topic, sorted by points
// descending.
func (store *ScoreStore) GetScoresByTopic(topicID int) ([]x.Score, error) {
	var scores []x.Score

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

	// Execute prepared statement
	if err := store.Select(&scores, query, topicID); err != nil {
		return []x.Score{}, fmt.Errorf("error getting scores: %w", err)
	}

	return scores, nil
}

// GetScoresByUser gets scores of a certain user, sorted by points descending.
func (store *ScoreStore) GetScoresByUser(userID int) ([]x.Score, error) {
	var scores []x.Score

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

	// Execute prepared statement
	if err := store.Select(&scores, query, userID); err != nil {
		return []x.Score{}, fmt.Errorf("error getting scores: %w", err)
	}

	return scores, nil
}

// GetScoresByTopicAndUser gets scores of a certain topic and user, sorted by
// points descending.
func (store *ScoreStore) GetScoresByTopicAndUser(topicID int, userID int) ([]x.Score, error) {
	var scores []x.Score

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

	// Execute prepared statement
	if err := store.Select(&scores, query, topicID, userID); err != nil {
		return []x.Score{}, fmt.Errorf("error getting scores: %w", err)
	}

	return scores, nil
}

// CreateScore creates a new score.
func (store *ScoreStore) CreateScore(score *x.Score) error {

	query := `
		INSERT INTO scores(topic_id, user_id, points, date) 
		VALUES (?, ?, ?, ?)
		`

	// Execute prepared statement
	if _, err := store.Exec(query,
		score.TopicID,
		score.UserID,
		score.Points,
		time.Now(),
	); err != nil {
		return fmt.Errorf("error creating score: %w", err)
	}

	return nil
}
