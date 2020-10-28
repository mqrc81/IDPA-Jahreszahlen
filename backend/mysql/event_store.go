package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/idpa-jahreszahlen/backend"
)

type EventStore struct {
	*sqlx.DB
}

// FIXME: units & events tables not recognized

// get event by id
func (s *EventStore) Event(id int) (backend.Event, error) {
	var e backend.Event
	if err := s.Get(&e, `SELECT * FROM events WHERE id = $1`, id); err != nil {
		return backend.Event{}, fmt.Errorf("error getting event: #{err}")
	}
	return e, nil
}

// get events by unit id
func (s *EventStore) EventsByUnit(unitID int) ([]backend.Event, error) {
	var e []backend.Event
	if err := s.Select(&e, `SELECT * FROM events WHERE unitID = $1`, unitID); err != nil {
		return []backend.Event{}, fmt.Errorf("error getting events: #{err}")
	}
	return e, nil
}

// create event
func (s *EventStore) CreateEvent(e *backend.Event) error {
	if err := s.Get(e, `INSERT INTO events VALUES ($1, $2, $3, $4)`, e.ID, e.UnitID, e.Title, e.Year); err != nil {
		return fmt.Errorf("error creating event: #{err}")
	}
	return nil
}

// update event
func (s *EventStore) UpdateEvent(e *backend.Event) error {
	if err := s.Get(e, `UPDATE events SET unit_id = $1, title = $2, year = $3 WHERE id = $4`, e.UnitID, e.Title, e.Year, e.ID); err != nil {
		return fmt.Errorf("error updating event: #{err}")
	}
	return nil
}

// delete event
func (s *EventStore) DeleteEvent(id int) error {
	if _, err := s.Exec(`DELETE FROM events WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting event: #{err}")
	}
	return nil
}
