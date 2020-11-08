package mysql

import (
	"fmt"
	//
	"github.com/jmoiron/sqlx"
	//
	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type EventStore struct {
	*sqlx.DB
}

/*
 * Get event by id
 */
func (s *EventStore) Event(id int) (backend.Event, error) {
	var e backend.Event
	if err := s.Get(&e, `SELECT * FROM events WHERE id = $1`, id); err != nil {
		return backend.Event{}, fmt.Errorf("error getting event: %w", err)
	}
	return e, nil
}

/*
 * Get events by unit id
 */
func (s *EventStore) EventsByUnit(unitID int) ([]backend.Event, error) {
	var ee []backend.Event
	if err := s.Select(&ee, `SELECT * FROM events WHERE unit_id = $1`, unitID); err != nil {
		return []backend.Event{}, fmt.Errorf("error getting events: %w", err)
	}
	return ee, nil
}

/*
 * Create event
 */
func (s *EventStore) CreateEvent(e *backend.Event) error {
	if err := s.Get(e, `INSERT INTO events VALUES ($1, $2, $3, $4)`, e.ID, e.UnitID, e.Title, e.Year); err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	return nil
}

/*
 * Update event
 */
func (s *EventStore) UpdateEvent(e *backend.Event) error {
	if err := s.Get(e, `UPDATE events SET unit_id = $1, title = $2, year = $3 WHERE id = $4`, e.UnitID, e.Title, e.Year, e.ID); err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}
	return nil
}

/*
 * Delete event
 */
func (s *EventStore) DeleteEvent(id int) error {
	if _, err := s.Exec(`DELETE FROM events WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}
	return nil
}
