package database

// event_store.go
// Part of the database layer. Contains all functions for events that access the
// database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// EventStore
// The MySQL database access object
type EventStore struct {
	*sqlx.DB
}

// Event
// Gets event by ID
func (store *EventStore) Event(eventID int) (backend.Event, error) {
	var event backend.Event

	// Execute prepared statement
	query := `
		SELECT e.*,
		       t.title AS topic_name
		FROM events e 
		    LEFT JOIN topics t ON t.topic_id = e.topic_id
		WHERE e.event_id = ?
		`
	if err := store.Get(&event, query, eventID); err != nil {
		return backend.Event{}, fmt.Errorf("error getting event: %w", err)
	}

	return event, nil
}

// EventsByTopic
// Gets events of a certain, sorted by year
func (store *EventStore) EventsByTopic(topicID int) ([]backend.Event, error) {
	var events []backend.Event

	// Execute prepared statement
	query := `
		SELECT e.*, 
		       t.title AS topic_name
		FROM events e 
		    LEFT JOIN topics t ON t.topic_id = e.topic_id
		WHERE e.topic_id = ?
		ORDER BY e.year
		`
	if err := store.Select(&events, query, topicID); err != nil {
		return []backend.Event{}, fmt.Errorf("error getting events: %w", err)
	}

	return events, nil
}

// CountEvents
// Gets amount of events
func (store *EventStore) CountEvents() (int, error) {
	var eventCount int

	// Execute prepared statement
	query := `
		SELECT COUNT(*) 
		FROM events
		`
	if err := store.Get(&eventCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of events: %w", err)
	}

	return eventCount, nil
}

// CreateEvent
// Creates a new event
func (store *EventStore) CreateEvent(event *backend.Event) error {

	// Execute prepared statement
	query := `
		INSERT INTO events(topic_id, title, year) 
		VALUES (?, ?, ?)
		`
	if _, err := store.Exec(query,
		event.TopicID,
		event.Title,
		event.Year); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}

	return nil
}

// UpdateEvent
// Updates an existing event
func (store *EventStore) UpdateEvent(event *backend.Event) error {

	// Execute prepared statement
	query := `
		UPDATE events 
		SET title = ?, 
		    year = ? 
		WHERE event_id = ?
		`
	if _, err := store.Exec(query,
		event.Title,
		event.Year,
		event.EventID); err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}

	return nil
}

// DeleteEvent
// Deletes an existing event
func (store *EventStore) DeleteEvent(eventID int) error {

	// Execute prepared statement
	query := `
		DELETE FROM events 
		WHERE event_id = ?
		`
	if _, err := store.Exec(query, eventID); err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}

	return nil
}
