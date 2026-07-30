package store

import (
	"database/sql"
	"errors"

	"coffee-site-server/internal/models"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("already exists")

func (s *Store) CreateUser(email, passwordHash string, permission models.Permission, verified bool) (*models.User, error) {
	res, err := s.DB.Exec(
		`INSERT INTO users (email, password_hash, permission, verified) VALUES (?, ?, ?, ?)`,
		email, passwordHash, permission, boolToInt(verified),
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetUserByID(id)
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	row := s.DB.QueryRow(
		`SELECT id, email, password_hash, permission, verified, created_at FROM users WHERE email = ?`,
		email,
	)
	return scanUser(row)
}

func (s *Store) GetUserByID(id int64) (*models.User, error) {
	row := s.DB.QueryRow(
		`SELECT id, email, password_hash, permission, verified, created_at FROM users WHERE id = ?`,
		id,
	)
	return scanUser(row)
}

func (s *Store) MarkUserVerified(email string) error {
	_, err := s.DB.Exec(`UPDATE users SET verified = 1 WHERE email = ?`, email)
	return err
}

func (s *Store) UpdateUserPassword(email, passwordHash string) error {
	res, err := s.DB.Exec(`UPDATE users SET password_hash = ? WHERE email = ?`, passwordHash, email)
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

func scanUser(row *sql.Row) (*models.User, error) {
	var u models.User
	var verified int
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Permission, &verified, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.Verified = verified != 0
	return &u, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
