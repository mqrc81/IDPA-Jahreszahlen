package database

// topic_store.go
// Part of the database layer. Contains all functions for topics that access the
// database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// TopicStore
// The database access object
type TopicStore struct {
	*sqlx.DB
}

// Topic
// Gets a topic by ID
func (store *TopicStore) Topic(topicID int) (backend.Topic, error) {
	var t backend.Topic

	// Execute prepared statement
	query := `SELECT * FROM topics WHERE topic_id = ?`
	if err := store.Get(&t, query, topicID); err != nil {
		return backend.Topic{}, fmt.Errorf("error getting topic: %w", err)
	}
	return t, nil
}

// Topics
// Gets all topics
func (store *TopicStore) Topics() ([]backend.Topic, error) {
	var tt []backend.Topic

	// Execute prepared statement
	query := `SELECT * FROM topics ORDER BY start_year`
	if err := store.Select(&tt, query); err != nil {
		return []backend.Topic{}, fmt.Errorf("error getting topics: %w", err)
	}
	return tt, nil
}

// CountTopics
// Gets amount of topics
func (store *EventStore) CountTopics() (int, error) {
	var tCount int

	// Execute prepared statement
	query := `SELECT COUNT(*) FROM topics`
	if err := store.Get(&tCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of topics: %w", err)
	}
	return tCount, nil
}

// CreateTopic
// Creates a new topic
func (store *TopicStore) CreateTopic(t *backend.Topic) error {

	// Execute prepared statement
	query := `INSERT INTO topics(title, start_year, end_year, description) VALUES (?, ?, ?, ?)`
	if _, err := store.Exec(query, t.Title,
		t.StartYear, t.EndYear, t.Description); err != nil {
		return fmt.Errorf("error creating t: %w", err)
	}

	// Execute prepared statement
	query = `SELECT * FROM topics WHERE topic_id = last_insert_id()`
	if err := store.Get(t, query); err != nil {
		return fmt.Errorf("error getting created topic: %w", err)
	}
	return nil
}

// UpdateTopic
// Updates an existing topic
func (store *TopicStore) UpdateTopic(t *backend.Topic) error {

	// Execute prepared statement
	query := `UPDATE topics SET title = ?, start_year = ?, end_year = ?, description = ? WHERE topic_id = ?`
	if _, err := store.Exec(query, t.Title, t.StartYear, t.EndYear, t.Description, t.TopicID); err != nil {
		return fmt.Errorf("error updating t: %w", err)
	}

	// Execute prepared statement
	query = `SELECT * FROM topics WHERE topic_id = last_insert_id()`
	if err := store.Get(t, query); err != nil {
		return fmt.Errorf("error getting updated topic: %w", err)
	}
	return nil
}

// DeleteTopic
// Deletes an existing topic
func (store *TopicStore) DeleteTopic(topicID int) error {

	// Execute prepared statement
	query := `DELETE FROM topics WHERE topic_id = ?`
	if _, err := store.Exec(query, topicID); err != nil {
		return fmt.Errorf("error deleting topic: %w", err)
	}
	return nil
}
