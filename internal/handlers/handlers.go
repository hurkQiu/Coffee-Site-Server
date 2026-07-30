package handlers

import (
	"coffee-site-server/internal/config"
	"coffee-site-server/internal/store"
)

type Handler struct {
	Store *store.Store
	Cfg   config.Config
}

func New(s *store.Store, cfg config.Config) *Handler {
	return &Handler{Store: s, Cfg: cfg}
}
