package mysql

/*
 * score_store.go contains all functions for scores that require database access
 */

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type ScoreStore struct {
	*sqlx.DB
}

/*
 * Scores gets a certain amount of scores sorted by points
 */
func (store ScoreStore) Scores(limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := store.Select(&ss, `SELECT * FROM scores ORDER BY points DESC LIMIT ?, ?`,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * ScoresByTopic gets a certain amount of scores by topic ID sorted by points
 */
func (store ScoreStore) ScoresByTopic(topicID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := store.Select(&ss, `SELECT * FROM scores WHERE topic_id = ? ORDER BY points DESC LIMIT ?, ?`,
		topicID,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * ScoresByUser gets a certain amount of scores by user ID sorted by points
 */
func (store ScoreStore) ScoresByUser(userID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := store.Select(&ss, `SELECT * FROM scores WHERE user_id = ? ORDER BY points DESC LIMIT ?, ?`,
		userID,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * ScoresByTopicAndUser gets a certain amount of scores by topic ID and user ID sorted by points
 */
func (store ScoreStore) ScoresByTopicAndUser(topicID int, userID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score
	query := `SELECT * FROM scores WHERE topic_id = ? AND user_id = ? ORDER BY points DESC LIMIT ?, ?`
	if err := store.Select(&ss, query, topicID, userID, offset, limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * CreateScore creates score
 */
func (store ScoreStore) CreateScore(s *backend.Score) error {
	query := `INSERT INTO scores(topic_id, user_id, points, date) VALUES (?, ?, ?, ?)`
	if _, err := store.Exec(query,
		s.TopicID,
		s.UserID,
		s.Points,
		s.Date); err != nil {
		return fmt.Errorf("error creating s: %w", err)
	}
	if err := store.Get(s, `SELECT * FROM scores WHERE score_id = last_insert_id()`); err != nil {
		return fmt.Errorf("error getting created s: %w", err)
	}
	return nil
}
