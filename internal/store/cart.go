package store

import (
	"coffee-site-server/internal/models"
)

func (s *Store) ListCartItems(guestID string) ([]models.CartItem, error) {
	rows, err := s.DB.Query(
		`SELECT product_id, name, price, image, quantity FROM cart_items WHERE guest_id = ? ORDER BY id`,
		guestID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.CartItem{}
	for rows.Next() {
		var it models.CartItem
		if err := rows.Scan(&it.ID, &it.Name, &it.Price, &it.Image, &it.Quantity); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (s *Store) AddOrIncrementCartItem(guestID string, item models.CartItem) error {
	_, err := s.DB.Exec(
		`INSERT INTO cart_items (guest_id, product_id, name, price, image, quantity)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT (guest_id, product_id) DO UPDATE SET quantity = quantity + excluded.quantity`,
		guestID, item.ID, item.Name, item.Price, item.Image, item.Quantity,
	)
	return err
}

func (s *Store) SetCartItemQuantity(guestID, productID string, quantity int) error {
	if quantity <= 0 {
		return s.RemoveCartItem(guestID, productID)
	}
	_, err := s.DB.Exec(
		`UPDATE cart_items SET quantity = ? WHERE guest_id = ? AND product_id = ?`,
		quantity, guestID, productID,
	)
	return err
}

func (s *Store) RemoveCartItem(guestID, productID string) error {
	_, err := s.DB.Exec(`DELETE FROM cart_items WHERE guest_id = ? AND product_id = ?`, guestID, productID)
	return err
}

func (s *Store) ClearCart(guestID string) error {
	_, err := s.DB.Exec(`DELETE FROM cart_items WHERE guest_id = ?`, guestID)
	return err
}
