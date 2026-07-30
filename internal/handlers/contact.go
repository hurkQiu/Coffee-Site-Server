package handlers

import (
	"net/http"

	"coffee-site-server/internal/httpx"
	"coffee-site-server/internal/middleware"
)

type contactRequest struct {
	OrderID string `json:"orderId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *Handler) CreateContactMessage(w http.ResponseWriter, r *http.Request) {
	var req contactRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.Title == "" || req.Content == "" {
		httpx.Error(w, http.StatusBadRequest, "請填寫標題與內容")
		return
	}

	var userID *int64
	if claims, ok := middleware.ClaimsFromContext(r.Context()); ok {
		userID = &claims.UserID
	}

	msg, err := h.Store.CreateContactMessage(userID, req.OrderID, req.Title, req.Content)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "送出失敗，請稍後再試")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, msg)
}
