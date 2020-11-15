package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type EventStore struct {
	*sqlx.DB
}

/*
 * Get event by event id
 */
func (s *EventStore) Event(eventID int) (backend.Event, error) {
	var e backend.Event
	if err := s.Get(&e, `SELECT * FROM events WHERE event_id = $1`, eventID); err != nil {
		return backend.Event{}, fmt.Errorf("error getting event: %w", err)
	}
	return e, nil
}

/*
 * Get events by topic id
 */
func (s *EventStore) EventsByTopic(topicID int) ([]backend.Event, error) {
	var ee []backend.Event
	if err := s.Select(&ee, `SELECT * FROM events WHERE topic_id = $1`, topicID); err != nil {
		return []backend.Event{}, fmt.Errorf("error getting events: %w", err)
	}
	return ee, nil
}

/*
 * Create event
 */
func (s *EventStore) CreateEvent(e *backend.Event) error {
	if _, err := s.Exec(`INSERT INTO events(topic_id, title, year) VALUES ($1, $2, $3)`,
		e.TopicID,
		e.Title,
		e.Year); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	return nil
}

/*
 * Update event
 */
func (s *EventStore) UpdateEvent(e *backend.Event) error {
	if _, err := s.Exec(`UPDATE events SET title = $1, year = $2 WHERE event_id = $2`,
		e.Title,
		e.Year,
		e.EventID); err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}
	return nil
}

/*
 * Delete event
 */
func (s *EventStore) DeleteEvent(eventID int) error {
	if _, err := s.Exec(`DELETE FROM events WHERE event_id = $1`, eventID); err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}
	return nil
}
