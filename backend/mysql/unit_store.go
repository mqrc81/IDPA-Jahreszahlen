package mysql

import (
	"fmt"
	//
	"github.com/jmoiron/sqlx"
	//
	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type UnitStore struct {
	*sqlx.DB
}

/*
 * Get unit by id
 */
func (s *UnitStore) Unit(id int) (backend.Unit, error) {
	var u backend.Unit
	if err := s.Get(&u, `SELECT * FROM units WHERE id = $1`, id); err != nil {
		return backend.Unit{}, fmt.Errorf("error getting unit: %w", err)
	}
	return u, nil
}

/*
 * Get units by id
 */
func (s *UnitStore) Units() ([]backend.Unit, error) {
	var uu []backend.Unit
	if err := s.Select(&uu, `SELECT * FROM units`); err != nil {
		return []backend.Unit{}, fmt.Errorf("error getting units: %w", err)
	}
	return uu, nil
}

/*
 * Create unit
 */
func (s *UnitStore) CreateUnit(u *backend.Unit) error {
	if err := s.Get(&u, `INSERT INTO units VALUES ($1, $2, $3, $4, $5, $6)`,
		u.ID, u.Title, u.StartYear, u.EndYear, u.Description, u.PlayCount); err != nil {
		return fmt.Errorf("error creating unit: %w", err)
	}
	return nil
}

/*
 * Update unit
 */
func (s *UnitStore) UpdateUnit(u *backend.Unit) error {
	if err := s.Get(&u, `UPDATE units SET title = $1, start_year = $2, end_year = $3, description = $4 WHERE id = $5`,
		u.Title, u.StartYear, u.EndYear, u.Description, u.ID); err != nil {
		return fmt.Errorf("error updating unit: %w", err)
	}
	return nil
}

/*
 * Delete unit by id
 */
func (s *UnitStore) DeleteUnit(id int) error {
	if _, err := s.Exec(`DELETE FROM units WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting unit: %w", err)
	}
	return nil
}
