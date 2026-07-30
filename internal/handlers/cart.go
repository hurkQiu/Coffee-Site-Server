package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"coffee-site-server/internal/httpx"
	"coffee-site-server/internal/models"
)

func guestIDFromRequest(r *http.Request) (string, bool) {
	id := r.Header.Get("X-Guest-Id")
	return id, id != ""
}

func (h *Handler) ListCart(w http.ResponseWriter, r *http.Request) {
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
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	guestID, ok := guestIDFromRequest(r)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "缺少 X-Guest-Id")
		return
	}
	var item models.CartItem
	if err := httpx.DecodeJSON(r, &item); err != nil || item.ID == "" || item.Quantity <= 0 {
		httpx.Error(w, http.StatusBadRequest, "無效的購物車資料")
		return
	}
	if err := h.Store.AddOrIncrementCartItem(guestID, item); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "加入購物車失敗")
		return
	}
	items, err := h.Store.ListCartItems(guestID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

type updateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

func (h *Handler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	guestID, ok := guestIDFromRequest(r)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "缺少 X-Guest-Id")
		return
	}
	productID := chi.URLParam(r, "productId")
	var req updateCartItemRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的請求內容")
		return
	}
	if err := h.Store.SetCartItemQuantity(guestID, productID, req.Quantity); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "更新購物車失敗")
		return
	}
	items, err := h.Store.ListCartItems(guestID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	guestID, ok := guestIDFromRequest(r)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "缺少 X-Guest-Id")
		return
	}
	productID := chi.URLParam(r, "productId")
	if err := h.Store.RemoveCartItem(guestID, productID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "移除商品失敗")
		return
	}
	items, err := h.Store.ListCartItems(guestID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	guestID, ok := guestIDFromRequest(r)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, "缺少 X-Guest-Id")
		return
	}
	if err := h.Store.ClearCart(guestID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "清空購物車失敗")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, []models.CartItem{})
}
