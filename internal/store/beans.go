package store

import (
	"database/sql"
	"errors"

	"coffee-site-server/internal/models"
)

func (s *Store) ListBeans() ([]models.CoffeeBean, error) {
	rows, err := s.DB.Query(`SELECT id, name, price, image, stock, roast, process FROM coffee_beans ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beans := []models.CoffeeBean{}
	for rows.Next() {
		var b models.CoffeeBean
		if err := rows.Scan(&b.ID, &b.Name, &b.Price, &b.Image, &b.Stock, &b.Roast, &b.Process); err != nil {
			return nil, err
		}
		beans = append(beans, b)
	}
	return beans, rows.Err()
}

func (s *Store) GetBean(id int64) (*models.CoffeeBean, error) {
	row := s.DB.QueryRow(`SELECT id, name, price, image, stock, roast, process FROM coffee_beans WHERE id = ?`, id)
	var b models.CoffeeBean
	if err := row.Scan(&b.ID, &b.Name, &b.Price, &b.Image, &b.Stock, &b.Roast, &b.Process); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (s *Store) CreateBean(b models.CoffeeBean) (*models.CoffeeBean, error) {
	res, err := s.DB.Exec(
		`INSERT INTO coffee_beans (name, price, image, stock, roast, process) VALUES (?, ?, ?, ?, ?, ?)`,
		b.Name, b.Price, b.Image, b.Stock, b.Roast, b.Process,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetBean(id)
}

func (s *Store) UpdateBean(id int64, b models.CoffeeBean) (*models.CoffeeBean, error) {
	res, err := s.DB.Exec(
		`UPDATE coffee_beans SET name = ?, price = ?, image = ?, stock = ?, roast = ?, process = ? WHERE id = ?`,
		b.Name, b.Price, b.Image, b.Stock, b.Roast, b.Process, id,
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
	return s.GetBean(id)
}

func (s *Store) DeleteBean(id int64) error {
	res, err := s.DB.Exec(`DELETE FROM coffee_beans WHERE id = ?`, id)
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
