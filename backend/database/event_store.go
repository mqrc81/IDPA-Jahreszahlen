package database

/*
 * event_store.go contains all functions for events that require database access
 */

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

/*
 * EventStore implements database access
 */
type EventStore struct {
	*sqlx.DB
}

/*
 * Event gets event by event ID
 */
func (store *EventStore) Event(eventID int) (backend.Event, error) {
	var e backend.Event
	query := `SELECT * FROM events WHERE event_id = ?`
	if err := store.Get(&e, query, eventID); err != nil {
		return backend.Event{}, fmt.Errorf("error getting event: %w", err)
	}
	return e, nil
}

/*
 * EventsByTopic gets events by topic ID, sorted randomly or by year
 */
func (store *EventStore) EventsByTopic(topicID int, orderByRand bool) ([]backend.Event, error) {
	var ee []backend.Event
	query := `SELECT * FROM events WHERE topic_id = ? ORDER BY year`
	if orderByRand {
		query = `SELECT * FROM events WHERE topic_id = ? ORDER BY RAND()`
	}
	if err := store.Select(&ee, query, topicID); err != nil {
		return []backend.Event{}, fmt.Errorf("error getting events: %w", err)
	}
	return ee, nil
}

/*
 * EventsCount gets number of events
 */
func (store *EventStore) EventsCount() (int, error) {
	var eCount int
	query := `SELECT COUNT(*) FROM events`
	if err := store.Get(&eCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of events: %w", err)
	}
	return eCount, nil
}

/*
 * CreateEvent creates event
 */
func (store *EventStore) CreateEvent(e *backend.Event) error {
	query := `INSERT INTO events(topic_id, title, year) VALUES (?, ?, ?)`
	if _, err := store.Exec(query, e.TopicID, e.Title, e.Year); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	query = `SELECT * FROM events WHERE event_id = last_insert_id()`
	if err := store.Get(e, query); err != nil {
		return fmt.Errorf("error getting created event: %w", err)
	}
	return nil
}

/*
 * UpdateEvent updates event
 */
func (store *EventStore) UpdateEvent(e *backend.Event) error {
	query := `UPDATE events SET title = ?, year = ? WHERE event_id = ?`
	if _, err := store.Exec(query, e.Title, e.Year, e.EventID); err != nil {
		return fmt.Errorf("error updating e: %w", err)
	}
	query = `SELECT * FROM events WHERE event_id = last_insert_id()`
	if err := store.Get(e, query); err != nil {
		return fmt.Errorf("error getting updated event: %w", err)
	}
	return nil
}

/*
 * DeleteEvent deletes event by event ID
 */
func (store *EventStore) DeleteEvent(eventID int) error {
	query := `DELETE FROM events WHERE event_id = ?`
	if _, err := store.Exec(query, eventID); err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}
	return nil
}
