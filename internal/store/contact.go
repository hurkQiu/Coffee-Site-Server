package store

import (
	"database/sql"

	"coffee-site-server/internal/models"
)

func (s *Store) CreateContactMessage(userID *int64, orderID, title, content string) (*models.ContactMessage, error) {
	var uid sql.NullInt64
	if userID != nil {
		uid = sql.NullInt64{Int64: *userID, Valid: true}
	}
	res, err := s.DB.Exec(
		`INSERT INTO contact_messages (user_id, order_id, title, content) VALUES (?, ?, ?, ?)`,
		uid, orderID, title, content,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	row := s.DB.QueryRow(
		`SELECT id, user_id, order_id, title, content, created_at FROM contact_messages WHERE id = ?`,
		id,
	)
	var msg models.ContactMessage
	var scannedUID sql.NullInt64
	if err := row.Scan(&msg.ID, &scannedUID, &msg.OrderID, &msg.Title, &msg.Content, &msg.CreatedAt); err != nil {
		return nil, err
	}
	if scannedUID.Valid {
		msg.UserID = &scannedUID.Int64
	}
	return &msg, nil
}
