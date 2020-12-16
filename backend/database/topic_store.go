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
	var topic backend.Topic

	// Execute prepared statement
	query := `SELECT * FROM topics WHERE topic_id = ?`
	if err := store.Get(&topic, query, topicID); err != nil {
		return backend.Topic{}, fmt.Errorf("error getting topic: %w", err)
	}
	return topic, nil
}

// Topics
// Gets all topics
func (store *TopicStore) Topics() ([]backend.Topic, error) {
	var topics []backend.Topic

	// Execute prepared statement
	query := `SELECT * FROM topics ORDER BY start_year`
	if err := store.Select(&topics, query); err != nil {
		return []backend.Topic{}, fmt.Errorf("error getting topics: %w", err)
	}
	return topics, nil
}

// CountTopics
// Gets amount of topics
func (store *EventStore) CountTopics() (int, error) {
	var topicCount int

	// Execute prepared statement
	query := `SELECT COUNT(*) FROM topics`
	if err := store.Get(&topicCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of topics: %w", err)
	}
	return topicCount, nil
}

// CreateTopic
// Creates a new topic
func (store *TopicStore) CreateTopic(topic *backend.Topic) error {

	// Execute prepared statement
	query := `INSERT INTO topics(title, start_year, end_year, description) VALUES (?, ?, ?, ?)`
	if _, err := store.Exec(query,
		topic.Title,
		topic.StartYear,
		topic.EndYear,
		topic.Description); err != nil {
		return fmt.Errorf("error creating topic: %w", err)
	}

	// Execute prepared statement
	query = `SELECT * FROM topics WHERE topic_id = last_insert_id()`
	if err := store.Get(topic, query); err != nil {
		return fmt.Errorf("error getting created topic: %w", err)
	}
	return nil
}

// UpdateTopic
// Updates an existing topic
func (store *TopicStore) UpdateTopic(topic *backend.Topic) error {

	// Execute prepared statement
	query := `UPDATE topics SET title = ?, start_year = ?, end_year = ?, description = ? WHERE topic_id = ?`
	if _, err := store.Exec(query,
		topic.Title,
		topic.StartYear,
		topic.EndYear,
		topic.Description,
		topic.TopicID); err != nil {
		return fmt.Errorf("error updating topic: %w", err)
	}

	// Execute prepared statement
	query = `SELECT * FROM topics WHERE topic_id = last_insert_id()`
	if err := store.Get(topic, query); err != nil {
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
