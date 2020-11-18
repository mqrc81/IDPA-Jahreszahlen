package mysql

/*
 * TODO Header
 */

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type EventStore struct {
	*sqlx.DB
}

/*
 * Event gets event by event ID
 */
func (s *EventStore) Event(eventID int) (backend.Event, error) {
	var e backend.Event
	if err := s.Get(&e, `SELECT * FROM events WHERE event_id = ?`, eventID); err != nil {
		return backend.Event{}, fmt.Errorf("error getting event: %w", err)
	}
	return e, nil
}

/*
 * EventsByTopic gets events by topic ID, sorted randomly or by year
 */
func (s *EventStore) EventsByTopic(topicID int, orderByRand bool) ([]backend.Event, error) {
	var ee []backend.Event
	order := "year"
	if orderByRand {
		order = "RAND()"
	}
	if err := s.Select(&ee, `SELECT * FROM events WHERE topic_id = ? ORDER BY ?`, topicID, order); err != nil {
		return []backend.Event{}, fmt.Errorf("error getting events: %w", err)
	}
	return ee, nil
}

/*
 * CreateEvent creates event
 */
func (s *EventStore) CreateEvent(event *backend.Event) error {
	if _, err := s.Exec(`INSERT INTO events(topic_id, title, year) VALUES (?, ?, ?)`,
		event.TopicID,
		event.Title,
		event.Year); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	if err := s.Get(event, `SELECT * FROM events WHERE event_id = last_insert_id()`); err != nil {
		return fmt.Errorf("error getting created event: %w", err)
	}
	return nil
}

/*
 * UpdateEvent updates event
 */
func (s *EventStore) UpdateEvent(event *backend.Event) error {
	if _, err := s.Exec(`UPDATE events SET title = ?, year = ? WHERE event_id = ?`,
		event.Title,
		event.Year,
		event.EventID); err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}
	if err := s.Get(event, `SELECT * FROM events WHERE event_id = last_insert_id()`); err != nil {
		return fmt.Errorf("error getting updated event: %w", err)
	}
	return nil
}

/*
 * DeleteEvent deletes event by event ID
 */
func (s *EventStore) DeleteEvent(eventID int) error {
	if _, err := s.Exec(`DELETE FROM events WHERE event_id = ?`, eventID); err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}
	return nil
}
