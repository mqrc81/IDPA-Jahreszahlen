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
func (s ScoreStore) Scores(limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := s.Select(&ss, `SELECT * FROM scores ORDER BY points DESC LIMIT ?, ?`,
		offset,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * ScoresByTopic gets a certain amount of scores by topic ID sorted by points
 */
func (s ScoreStore) ScoresByTopic(topicID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := s.Select(&ss, `SELECT * FROM scores WHERE topic_id = ? ORDER BY points DESC LIMIT ?, ?`,
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
func (s ScoreStore) ScoresByUser(userID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := s.Select(&ss, `SELECT * FROM scores WHERE user_id = ? ORDER BY points DESC LIMIT ?, ?`,
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
func (s ScoreStore) ScoresByTopicAndUser(topicID int, userID int, limit int, offset int) ([]backend.Score, error) {
	var ss []backend.Score
	query := `SELECT * FROM scores WHERE topic_id = ? AND user_id = ? ORDER BY points DESC LIMIT ?, ?`
	if err := s.Select(&ss, query, topicID, userID, offset, limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * CreateScore creates score
 */
func (s ScoreStore) CreateScore(score *backend.Score) error {
	if _, err := s.Exec(`INSERT INTO scores(topic_id, user_id, points, date) VALUES (?, ?, ?, ?)`,
		score.TopicID,
		score.UserID,
		score.Points,
		score.Date); err != nil {
		return fmt.Errorf("error creating score: %w", err)
	}
	if err := s.Get(score, `SELECT * FROM scores WHERE score_id = last_insert_id()`); err != nil {
	    return fmt.Errorf("error getting created score: %w", err)
	}
	return nil
}
