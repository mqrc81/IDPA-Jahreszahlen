package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mqrc81/idpa-jahreszahlen/backend"
)

type UnitStore struct {
	*sqlx.DB
}

// FIXME: units & events tables not recognized

// get unit by id
func (s *UnitStore) Unit(id int) (backend.Unit, error) {
	var u backend.Unit
	if err := s.Get(&u, `SELECT * FROM units WHERE id = $1`, id); err != nil {
		return backend.Unit{}, fmt.Errorf("error getting unit: #{err}")
	}
	return u, nil
}

// get units by id
func (s *UnitStore) Units() ([]backend.Unit, error) {
	var u []backend.Unit
	if err := s.Get(&u, `SELECT * FROM units`); err != nil {
		return []backend.Unit{}, fmt.Errorf("error getting units: #{err}")
	}
	return u, nil
}

// create unit
func (s *UnitStore) CreateUnit(u *backend.Unit) error {
	if err := s.Get(&u, `INSERT INTO units VALUES ($1, $2, $3, $4)`, u.ID, u.Title, u.Description, u.PlayCount); err != nil {
		return fmt.Errorf("error creating unit: #{err}")
	}
	return nil
}

// update unit
func (s *UnitStore) UpdateUnit(u *backend.Unit) error {
	if err := s.Get(&u, `UPDATE units SET title = $1, description = $2 WHERE id = $3`, u.Title, u.Description, u.ID); err != nil {
		return fmt.Errorf("error updating unit: #{err}")
	}
	return nil
}

// delete unit by id
func (s *UnitStore) DeleteUnit(id int) error {
	if _, err := s.Exec(`DELETE FROM units WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting unit: #{err}")
	}
	return nil
}
