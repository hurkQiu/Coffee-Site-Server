package store

import (
	"fmt"
	"math/rand"
	"time"

	"coffee-site-server/internal/models"
)

func (s *Store) CreateOrder(userID int64, items []models.CartItem) (*models.Order, error) {
	total := 0
	for _, it := range items {
		total += it.Price * it.Quantity
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	orderNumber := generateOrderNumber()
	res, err := tx.Exec(
		`INSERT INTO orders (order_number, user_id, total_price) VALUES (?, ?, ?)`,
		orderNumber, userID, total,
	)
	if err != nil {
		return nil, err
	}
	orderID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	for _, it := range items {
		if _, err := tx.Exec(
			`INSERT INTO order_items (order_id, name, image, quantity, price) VALUES (?, ?, ?, ?, ?)`,
			orderID, it.Name, it.Image, it.Quantity, it.Price,
		); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetOrder(orderID)
}

func (s *Store) GetOrder(id int64) (*models.Order, error) {
	row := s.DB.QueryRow(`SELECT id, order_number, user_id, total_price, created_at FROM orders WHERE id = ?`, id)
	var o models.Order
	if err := row.Scan(&o.ID, &o.OrderNumber, &o.UserID, &o.TotalPrice, &o.CreatedAt); err != nil {
		return nil, err
	}
	items, err := s.listOrderItems(id)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}

func (s *Store) ListOrdersByUser(userID int64) ([]models.Order, error) {
	rows, err := s.DB.Query(
		`SELECT id, order_number, user_id, total_price, created_at FROM orders WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.OrderNumber, &o.UserID, &o.TotalPrice, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range orders {
		items, err := s.listOrderItems(orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}
	return orders, nil
}

func (s *Store) listOrderItems(orderID int64) ([]models.OrderItem, error) {
	rows, err := s.DB.Query(
		`SELECT name, image, quantity, price FROM order_items WHERE order_id = ? ORDER BY id`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.OrderItem{}
	for rows.Next() {
		var it models.OrderItem
		if err := rows.Scan(&it.Name, &it.Image, &it.Quantity, &it.Price); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func generateOrderNumber() string {
	return fmt.Sprintf("OD%s%03d", time.Now().Format("20060102150405"), rand.Intn(1000))
}
