package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	_ "modernc.org/sqlite"

	"coffee-site-server/internal/auth"
)

//go:embed schema.sql
var schemaSQL string

func Open(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", path)
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// modernc.org/sqlite does not support concurrent writers well; keep a single connection.
	conn.SetMaxOpenConns(1)

	if _, err := conn.Exec(schemaSQL); err != nil {
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	if err := seed(conn); err != nil {
		return nil, fmt.Errorf("seed: %w", err)
	}

	return conn, nil
}

func seed(conn *sql.DB) error {
	if err := seedOptions(conn, "roast_levels", []string{"淺焙", "中焙", "深焙"}); err != nil {
		return err
	}
	if err := seedOptions(conn, "process_methods", []string{"水洗", "日曬", "蜜處理", "厭氧"}); err != nil {
		return err
	}
	if err := seedOptions(conn, "utensil_categories", []string{"濾杯", "磨豆機", "耗材", "其餘用具"}); err != nil {
		return err
	}

	var beanCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM coffee_beans`).Scan(&beanCount); err != nil {
		return err
	}
	if beanCount == 0 {
		beans := []struct {
			name, roast, process string
			price, stock         int
		}{
			{"耶加雪菲", "淺焙", "水洗", 380, 12},
			{"西達摩", "淺焙", "日曬", 360, 8},
			{"藝伎莊園", "淺焙", "蜜處理", 520, 3},
			{"瑰夏厭氧", "淺焙", "厭氧", 620, 0},
			{"曼特寧", "中焙", "水洗", 340, 20},
			{"哥倫比亞", "中焙", "蜜處理", 350, 15},
			{"巴西聖多斯", "深焙", "日曬", 300, 25},
			{"摩卡爪哇", "深焙", "厭氧", 390, 6},
		}
		for _, b := range beans {
			if _, err := conn.Exec(
				`INSERT INTO coffee_beans (name, price, image, stock, roast, process) VALUES (?, ?, '', ?, ?, ?)`,
				b.name, b.price, b.stock, b.roast, b.process,
			); err != nil {
				return err
			}
		}
	}

	var utensilCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM utensils`).Scan(&utensilCount); err != nil {
		return err
	}
	if utensilCount == 0 {
		utensils := []struct {
			name, category string
			price, stock   int
		}{
			{"V60 濾杯", "濾杯", 450, 14},
			{"蛋糕型濾杯", "濾杯", 380, 9},
			{"手搖磨豆機", "磨豆機", 1200, 5},
			{"電動磨豆機", "磨豆機", 3200, 0},
			{"濾紙 100 入", "耗材", 150, 40},
			{"濾布", "耗材", 220, 18},
			{"手沖壺", "其餘用具", 890, 7},
			{"電子秤", "其餘用具", 650, 3},
		}
		for _, u := range utensils {
			if _, err := conn.Exec(
				`INSERT INTO utensils (name, price, image, stock, category) VALUES (?, ?, '', ?, ?)`,
				u.name, u.price, u.stock, u.category,
			); err != nil {
				return err
			}
		}
	}

	var userCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&userCount); err != nil {
		return err
	}
	if userCount == 0 {
		adminHash, err := auth.HashPassword("admin123")
		if err != nil {
			return err
		}
		memberHash, err := auth.HashPassword("member123")
		if err != nil {
			return err
		}
		if _, err := conn.Exec(
			`INSERT INTO users (email, password_hash, permission, verified) VALUES (?, ?, 'admin', 1)`,
			"admin@coffeehouse.example.com", adminHash,
		); err != nil {
			return err
		}
		if _, err := conn.Exec(
			`INSERT INTO users (email, password_hash, permission, verified) VALUES (?, ?, 'member', 1)`,
			"member@coffeehouse.example.com", memberHash,
		); err != nil {
			return err
		}
		log.Println("[seed] created default accounts admin@coffeehouse.example.com/admin123 and member@coffeehouse.example.com/member123")
	}

	return nil
}

func seedOptions(conn *sql.DB, table string, names []string) error {
	var count int
	if err := conn.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table)).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	for i, name := range names {
		if _, err := conn.Exec(
			fmt.Sprintf(`INSERT INTO %s (name, hidden, sort_order) VALUES (?, 0, ?)`, table),
			name, i,
		); err != nil {
			return err
		}
	}
	return nil
}
