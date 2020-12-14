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
	var e backend.Event

	// Execute prepared statement
	query := `SELECT * FROM events WHERE event_id = ?`
	if err := store.Get(&e, query, eventID); err != nil {
		return backend.Event{}, fmt.Errorf("error getting event: %w", err)
	}
	return e, nil
}

// EventsByTopic
// Gets events of a certain, sorted randomly or by year
func (store *EventStore) EventsByTopic(topicID int, orderByRand bool) ([]backend.Event, error) {
	var ee []backend.Event

	// Execute prepared statement
	query := `SELECT * FROM events WHERE topic_id = ? ORDER BY year`
	if orderByRand {
		query = `SELECT * FROM events WHERE topic_id = ? ORDER BY RAND()`
	}
	if err := store.Select(&ee, query, topicID); err != nil {
		return []backend.Event{}, fmt.Errorf("error getting events: %w", err)
	}
	return ee, nil
}

// CountEvents
// Gets amount of events
func (store *EventStore) CountEvents() (int, error) {
	var eCount int

	// Execute prepared statement
	query := `SELECT COUNT(*) FROM events`
	if err := store.Get(&eCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of events: %w", err)
	}
	return eCount, nil
}

// CountEventsByTopic
// Gets amount of events of a certain topic
func (store *EventStore) CountEventsByTopic(topicID int) (int, error) {
	var eCount int

	// Execute prepared statement
	query := `SELECT COUNT(*) FROM events WHERE topic_id = ?`
	if err := store.Get(&eCount, query, topicID); err != nil {
		return 0, fmt.Errorf("error getting number of events: %w", err)
	}
	return eCount, nil
}

// CreateEvent
// Creates a new event
func (store *EventStore) CreateEvent(e *backend.Event) error {

	// Execute prepared statement
	query := `INSERT INTO events(topic_id, title, year) VALUES (?, ?, ?)`
	if _, err := store.Exec(query, e.TopicID, e.Title, e.Year); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}

	// Execute prepared statement
	query = `SELECT * FROM events WHERE event_id = last_insert_id()`
	if err := store.Get(e, query); err != nil {
		return fmt.Errorf("error getting created event: %w", err)
	}
	return nil
}

// UpdateEvent
// Updates an existing event
func (store *EventStore) UpdateEvent(e *backend.Event) error {

	// Execute prepared statement
	query := `UPDATE events SET title = ?, year = ? WHERE event_id = ?`
	if _, err := store.Exec(query, e.Title, e.Year, e.EventID); err != nil {
		return fmt.Errorf("error updating e: %w", err)
	}

	// Execute prepared statement
	query = `SELECT * FROM events WHERE event_id = last_insert_id()`
	if err := store.Get(e, query); err != nil {
		return fmt.Errorf("error getting updated event: %w", err)
	}
	return nil
}

// DeleteEvent
// Deletes an existing event
func (store *EventStore) DeleteEvent(eventID int) error {

	// Execute prepared statement
	query := `DELETE FROM events WHERE event_id = ?`
	if _, err := store.Exec(query, eventID); err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}
	return nil
}
