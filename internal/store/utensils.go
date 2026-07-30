package store

import (
	"database/sql"
	"errors"

	"coffee-site-server/internal/models"
)

func (s *Store) ListUtensils() ([]models.Utensil, error) {
	rows, err := s.DB.Query(`SELECT id, name, price, image, stock, category FROM utensils ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	utensils := []models.Utensil{}
	for rows.Next() {
		var u models.Utensil
		if err := rows.Scan(&u.ID, &u.Name, &u.Price, &u.Image, &u.Stock, &u.Category); err != nil {
			return nil, err
		}
		utensils = append(utensils, u)
	}
	return utensils, rows.Err()
}

func (s *Store) GetUtensil(id int64) (*models.Utensil, error) {
	row := s.DB.QueryRow(`SELECT id, name, price, image, stock, category FROM utensils WHERE id = ?`, id)
	var u models.Utensil
	if err := row.Scan(&u.ID, &u.Name, &u.Price, &u.Image, &u.Stock, &u.Category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (s *Store) CreateUtensil(u models.Utensil) (*models.Utensil, error) {
	res, err := s.DB.Exec(
		`INSERT INTO utensils (name, price, image, stock, category) VALUES (?, ?, ?, ?, ?)`,
		u.Name, u.Price, u.Image, u.Stock, u.Category,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetUtensil(id)
}

func (s *Store) UpdateUtensil(id int64, u models.Utensil) (*models.Utensil, error) {
	res, err := s.DB.Exec(
		`UPDATE utensils SET name = ?, price = ?, image = ?, stock = ?, category = ? WHERE id = ?`,
		u.Name, u.Price, u.Image, u.Stock, u.Category, id,
	)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrNotFound
	}
	return s.GetUtensil(id)
}

func (s *Store) DeleteUtensil(id int64) error {
	res, err := s.DB.Exec(`DELETE FROM utensils WHERE id = ?`, id)
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
