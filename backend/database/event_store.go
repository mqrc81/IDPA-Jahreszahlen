// The database store evolving around events, with all necessary methods that
// access the database.

package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// EventStore is the MySQL database access object.
type EventStore struct {
	*sqlx.DB
}

// GetEvent gets event by ID.
func (store *EventStore) GetEvent(eventID int) (x.Event, error) {
	var event x.Event

	query := `
		SELECT *
		FROM events
		WHERE event_id = ?
		`

	// Execute prepared statement
	if err := store.Get(&event, query, eventID); err != nil {
		return x.Event{}, fmt.Errorf("error getting event: %w", err)
	}

	return event, nil
}

// CountEvents gets amount of events.
func (store *EventStore) CountEvents() (int, error) {
	var eventCount int

	query := `
		SELECT COUNT(*) 
		FROM events
		`

	// Execute prepared statement
	if err := store.Get(&eventCount, query); err != nil {
		return 0, fmt.Errorf("error getting number of events: %w", err)
	}

	return eventCount, nil
}

// CreateEvent creates a new event.
func (store *EventStore) CreateEvent(event *x.Event) error {

	query := `
		INSERT INTO events(topic_id, name, year, date) 
		VALUES (?, ?, ?, ?)
		`

	// Execute prepared statement
	if _, err := store.Exec(query,
		event.TopicID,
		event.Name,
		event.Year,
		event.Date,
	); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}

	return nil
}

// UpdateEvent updates an existing event.
func (store *EventStore) UpdateEvent(event *x.Event) error {

	query := `
		UPDATE events 
		SET name = ?, 
		    year = ?,
		    date = ?
		WHERE event_id = ?
		`

	// Execute prepared statement
	if _, err := store.Exec(query,
		event.Name,
		event.Year,
		event.Date,
		event.EventID,
	); err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}

	return nil
}

// DeleteEvent deletes an existing event.
func (store *EventStore) DeleteEvent(eventID int) error {

	query := `
		DELETE FROM events 
		WHERE event_id = ?
		`

	// Execute prepared statement
	if _, err := store.Exec(query, eventID); err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}

	return nil
}
