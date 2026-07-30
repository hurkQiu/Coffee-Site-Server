package handlers

import (
	"net/http"

	"coffee-site-server/internal/httpx"
	"coffee-site-server/internal/middleware"
)

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "請先登入")
		return
	}
	guestID, ok := guestIDFromRequest(r)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "缺少 X-Guest-Id")
		return
	}

	items, err := h.Store.ListCartItems(guestID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	if len(items) == 0 {
		httpx.Error(w, http.StatusBadRequest, "購物車是空的")
		return
	}

	order, err := h.Store.CreateOrder(claims.UserID, items)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "結帳失敗，請稍後再試")
		return
	}

	if err := h.Store.ClearCart(guestID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, order)
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "請先登入")
		return
	}
	orders, err := h.Store.ListOrdersByUser(claims.UserID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, orders)
}
