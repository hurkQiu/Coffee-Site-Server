package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"coffee-site-server/internal/httpx"
)

func (h *Handler) ListWishlist(w http.ResponseWriter, r *http.Request) {
	guestID, ok := guestIDFromRequest(r)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "缺少 X-Guest-Id")
		return
	}
	ids, err := h.Store.ListWishlist(guestID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, ids)
}

func (h *Handler) ToggleWishlistItem(w http.ResponseWriter, r *http.Request) {
	guestID, ok := guestIDFromRequest(r)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "缺少 X-Guest-Id")
		return
	}
	productID := chi.URLParam(r, "productId")
	wishlisted, err := h.Store.ToggleWishlistItem(guestID, productID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"id": productID, "wishlisted": wishlisted})
}
