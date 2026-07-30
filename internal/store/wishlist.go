package store

func (s *Store) ListWishlist(guestID string) ([]string, error) {
	rows, err := s.DB.Query(`SELECT product_id FROM wishlist_items WHERE guest_id = ? ORDER BY id`, guestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ToggleWishlistItem adds the product if absent, removes it if present, and
// reports whether it is now wishlisted.
func (s *Store) ToggleWishlistItem(guestID, productID string) (bool, error) {
	var exists int
	err := s.DB.QueryRow(
		`SELECT COUNT(*) FROM wishlist_items WHERE guest_id = ? AND product_id = ?`,
		guestID, productID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists > 0 {
		if _, err := s.DB.Exec(
			`DELETE FROM wishlist_items WHERE guest_id = ? AND product_id = ?`,
			guestID, productID,
		); err != nil {
			return false, err
		}
		return false, nil
	}

	if _, err := s.DB.Exec(
		`INSERT INTO wishlist_items (guest_id, product_id) VALUES (?, ?)`,
		guestID, productID,
	); err != nil {
		return false, err
	}
	return true, nil
}
