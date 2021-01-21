// The database store evolving around topics, with all necessary methods that
// access the database.

package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// TopicStore is the database access object.
type TopicStore struct {
	*sqlx.DB
}

// GetTopic gets a topic and its events by ID.
func (store *TopicStore) GetTopic(topicID int) (x.Topic, error) {
	var topic x.Topic

	// Execute prepared statement
	query := `
		SELECT t.*, 
		       COUNT(DISTINCT s.score_id) AS scores_count,
		       COUNT(DISTINCT e.event_id) AS events_count
		FROM topics t 
			LEFT JOIN scores s ON s.topic_id = t.topic_id 
		    LEFT JOIN events e on t.topic_id = e.topic_id
		WHERE t.topic_id = ?
		`
	if err := store.Get(&topic, query, topicID); err != nil {
		return x.Topic{}, fmt.Errorf("error getting topic: %w", err)
	}

	// Execute prepared statement
	query = `
		SELECT * 
		FROM events 
		WHERE topic_id = ? 
		ORDER BY date
		`
	if err := store.Select(&topic.Events, query, topicID); err != nil {
		return x.Topic{}, fmt.Errorf("error getting events of topic: %w", err)
	}

	return topic, nil
}

// GetTopics gets all topics and their events.
func (store *TopicStore) GetTopics() ([]x.Topic, error) {
	var topics []x.Topic

	// Execute prepared statement
	query := `
		SELECT t.*, 
		       COUNT(DISTINCT s.score_id) AS scores_count,
		       COUNT(DISTINCT e.event_id) AS events_count
		FROM topics t 
			LEFT JOIN scores s ON s.topic_id = t.topic_id 
		    LEFT JOIN events e on t.topic_id = e.topic_id
		GROUP BY t.topic_id, t.start_year 
		ORDER BY t.start_year
		`
	if err := store.Select(&topics, query); err != nil {
		return []x.Topic{}, fmt.Errorf("error getting topics: %w", err)
	}

	// Loop through topics to get events
	for _, topic := range topics {
		// Execute prepared statement
		query = `
		SELECT * 
		FROM events 
		WHERE topic_id = ?
		`
		if err := store.Select(&topic.Events, query, topic.TopicID); err != nil {
			return []x.Topic{}, fmt.Errorf("error getting events of topics: %w", err)
		}
	}

	return topics, nil
}

// CountTopics gets amount of topics.
func (store *EventStore) CountTopics() (int, error) {
	var topicCount int

	// Execute prepared statement
	query := `
		SELECT COUNT(*) 
		FROM topics
		`
	if err := store.Get(&topicCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of topics: %w", err)
	}

	return topicCount, nil
}

// CreateTopic creates a new topic.
func (store *TopicStore) CreateTopic(topic *x.Topic) error {

	// Execute prepared statement
	query := `
		INSERT INTO topics(name, start_year, end_year, description) 
		VALUES (?, ?, ?, ?)
		`
	if _, err := store.Exec(query,
		topic.Name,
		topic.StartYear,
		topic.EndYear,
		topic.Description); err != nil {
		return fmt.Errorf("error creating topic: %w", err)
	}

	return nil
}

// UpdateTopic updates an existing topic.
func (store *TopicStore) UpdateTopic(topic *x.Topic) error {

	// Execute prepared statement
	query := `
		UPDATE topics 
		SET name = ?, 
		    start_year = ?, 
		    end_year = ?, 
		    description = ? 
		WHERE topic_id = ?
		`
	if _, err := store.Exec(query,
		topic.Name,
		topic.StartYear,
		topic.EndYear,
		topic.Description,
		topic.TopicID); err != nil {
		return fmt.Errorf("error updating topic: %w", err)
	}

	return nil
}

// DeleteTopic deletes an existing topic.
func (store *TopicStore) DeleteTopic(topicID int) error {

	// Execute prepared statement
	query := `
		DELETE FROM topics 
		WHERE topic_id = ?
		`
	if _, err := store.Exec(query, topicID); err != nil {
		return fmt.Errorf("error deleting topic: %w", err)
	}

	return nil
}
