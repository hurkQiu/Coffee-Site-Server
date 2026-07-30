package store

import (
	"database/sql"
	"fmt"
	"strings"

	"coffee-site-server/internal/models"
)

// options.go backs the three simple "named option" tables shared by beans
// (roast levels, process methods) and utensils (categories). Table names are
// always one of a fixed internal set, never user input, so string formatting
// them into the query is safe.

func (s *Store) ListOptions(table string) ([]models.NamedOption, error) {
	rows, err := s.DB.Query(fmt.Sprintf(`SELECT name, hidden FROM %s ORDER BY sort_order, id`, table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := []models.NamedOption{}
	for rows.Next() {
		var opt models.NamedOption
		var hidden int
		if err := rows.Scan(&opt.Name, &hidden); err != nil {
			return nil, err
		}
		opt.Hidden = hidden != 0
		options = append(options, opt)
	}
	return options, rows.Err()
}

func (s *Store) AddOption(table, name string) error {
	var maxOrder sql.NullInt64
	if err := s.DB.QueryRow(fmt.Sprintf(`SELECT MAX(sort_order) FROM %s`, table)).Scan(&maxOrder); err != nil {
		return err
	}
	_, err := s.DB.Exec(fmt.Sprintf(`INSERT INTO %s (name, hidden, sort_order) VALUES (?, 0, ?)`, table), name, maxOrder.Int64+1)
	if err != nil && isUniqueConstraintErr(err) {
		return ErrConflict
	}
	return err
}

func (s *Store) ToggleOptionHidden(table, name string) error {
	res, err := s.DB.Exec(fmt.Sprintf(`UPDATE %s SET hidden = 1 - hidden WHERE name = ?`, table), name)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func isUniqueConstraintErr(err error) bool {
	return err != nil && strings.Contains(strings.ToUpper(err.Error()), "UNIQUE CONSTRAINT")
}
