package mysql

/*
 * TODO Header
 */

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type TopicStore struct {
	*sqlx.DB
}

/*
 * Topic gets topic by topic ID
 */
func (s *TopicStore) Topic(topicID int) (backend.Topic, error) {
	var t backend.Topic
	if err := s.Get(&t, `SELECT * FROM topics WHERE topic_id = ?`, topicID); err != nil {
		return backend.Topic{}, fmt.Errorf("error getting topic: %w", err)
	}
	return t, nil
}

/*
 * Topic gets topics
 */
func (s *TopicStore) Topics() ([]backend.Topic, error) {
	var tt []backend.Topic
	if err := s.Select(&tt, `SELECT * FROM topics ORDER BY start_year`); err != nil {
		return []backend.Topic{}, fmt.Errorf("error getting topics: %w", err)
	}
	return tt, nil
}

/*
 * CreateTopic creates topic
 */
func (s *TopicStore) CreateTopic(topic *backend.Topic) error {
	if _, err := s.Exec(`INSERT INTO topics(title, start_year, end_year, description) VALUES (?, ?, ?, ?)`,
		topic.Title,
		topic.StartYear,
		topic.EndYear,
		topic.Description); err != nil {
		return fmt.Errorf("error creating topic: %w", err)
	}
	return nil
}

/*
 * UpdateTopic updates topic
 */
func (s *TopicStore) UpdateTopic(topic *backend.Topic) error {
	if _, err := s.Exec(`UPDATE topics SET title = ?, start_year = ?, end_year = ?, description = ? WHERE topic_id = ?`,
		topic.Title,
		topic.StartYear,
		topic.EndYear,
		topic.Description,
		topic.TopicID); err != nil {
		return fmt.Errorf("error updating topic: %w", err)
	}
	return nil
}

/*
 * DeleteTopic deletes topic by topic ID
 */
func (s *TopicStore) DeleteTopic(topicID int) error {
	if _, err := s.Exec(`DELETE FROM topics WHERE topic_id = ?`, topicID); err != nil {
		return fmt.Errorf("error deleting topic: %w", err)
	}
	return nil
}
