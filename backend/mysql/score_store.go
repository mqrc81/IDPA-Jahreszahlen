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
 * ScoresByTopic gets scores by topic ID sorted by points
 */
func (s ScoreStore) ScoresByTopic(topicID int) ([]backend.Score, error) {
	var ss []backend.Score
	if err := s.Select(&ss, `SELECT * FROM scores WHERE topic_id = ? ORDER BY points DESC`, topicID); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * ScoresByUser gets scores by username sorted by points
 */
func (s ScoreStore) ScoresByUser(username string) ([]backend.Score, error) {
	var ss []backend.Score
	if err := s.Select(&ss, `SELECT * FROM scores WHERE username = ? ORDER BY points DESC`, username); err != nil {
		return []backend.Score{}, fmt.Errorf("error getting scores: %w", err)
	}
	return ss, nil
}

/*
 * CreateScore creates score
 */
func (s ScoreStore) CreateScore(score *backend.Score) error {
	if _, err := s.Exec(`INSERT INTO scores(topic_id, username, points) VALUES (?, ?, ?)`,
		score.TopicID,
		score.Username,
		score.Points); err != nil {
		return fmt.Errorf("error creating score: %w", err)
	}
	return nil
}
