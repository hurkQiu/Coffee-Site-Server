package store

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"coffee-site-server/internal/models"
)

const codeTTL = 10 * time.Minute

func GenerateCode() (string, error) {
	max := 1000000
	n, err := randInt(max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n), nil
}

func randInt(max int) (int, error) {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return 0, err
	}
	n := int(buf[0])<<24 | int(buf[1])<<16 | int(buf[2])<<8 | int(buf[3])
	if n < 0 {
		n = -n
	}
	return n % max, nil
}

func (s *Store) CreateVerificationCode(email string, purpose models.VerificationPurpose) (string, error) {
	code, err := GenerateCode()
	if err != nil {
		return "", err
	}
	// Invalidate previous outstanding codes for the same purpose so only the
	// latest one is usable.
	if _, err := s.DB.Exec(
		`UPDATE verification_codes SET consumed = 1 WHERE email = ? AND purpose = ? AND consumed = 0`,
		email, purpose,
	); err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(codeTTL)
	if _, err := s.DB.Exec(
		`INSERT INTO verification_codes (email, code, purpose, expires_at) VALUES (?, ?, ?, ?)`,
		email, code, purpose, expiresAt,
	); err != nil {
		return "", err
	}
	return code, nil
}

// ValidateVerificationCode checks that the code is correct and unexpired without consuming it.
func (s *Store) ValidateVerificationCode(email, code string, purpose models.VerificationPurpose) error {
	var expiresAt time.Time
	row := s.DB.QueryRow(
		`SELECT expires_at FROM verification_codes
		 WHERE email = ? AND code = ? AND purpose = ? AND consumed = 0
		 ORDER BY id DESC LIMIT 1`,
		email, code, purpose,
	)
	if err := row.Scan(&expiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCode
		}
		return err
	}
	if time.Now().After(expiresAt) {
		return ErrCodeExpired
	}
	return nil
}

func (s *Store) ConsumeVerificationCode(email, code string, purpose models.VerificationPurpose) error {
	if err := s.ValidateVerificationCode(email, code, purpose); err != nil {
		return err
	}
	_, err := s.DB.Exec(
		`UPDATE verification_codes SET consumed = 1
		 WHERE email = ? AND code = ? AND purpose = ? AND consumed = 0`,
		email, code, purpose,
	)
	return err
}

var ErrInvalidCode = errors.New("invalid code")
var ErrCodeExpired = errors.New("code expired")
