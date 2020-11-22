package mysql

/*
 * TODO Header
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
func (s ScoreStore) Scores(limit int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := s.Select(&ss, `SELECT * FROM scores ORDER BY points DESC LIMIT ?`,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * ScoresByTopic gets a certain amount of scores by topic ID sorted by points
 */
func (s ScoreStore) ScoresByTopic(topicID int, limit int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := s.Select(&ss, `SELECT * FROM scores WHERE topic_id = ? ORDER BY points DESC LIMIT ?`,
		topicID,
		limit); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * ScoresByUser gets a certain amount of scores by username sorted by points
 */
func (s ScoreStore) ScoresByUser(userID int, limit int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := s.Select(&ss, `SELECT * FROM scores WHERE user_id = ? ORDER BY points DESC LIMIT ?`,
		userID,
		limit); err != nil {
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
