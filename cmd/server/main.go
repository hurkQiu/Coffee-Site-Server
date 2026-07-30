package main

import (
	"log"
	"net/http"

	"coffee-site-server/internal/config"
	"coffee-site-server/internal/db"
	"coffee-site-server/internal/handlers"
	"coffee-site-server/internal/store"
)

func main() {
	cfg := config.Load()

	conn, err := db.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer conn.Close()

	s := store.New(conn)
	h := handlers.New(s, cfg)
	router := handlers.NewRouter(h)

	log.Printf("coffee-site-server listening on :%s (db=%s, allowed origin=%v)", cfg.Port, cfg.DatabasePath, cfg.AllowOrigins)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
