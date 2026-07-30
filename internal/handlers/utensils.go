package handlers

import (
	"net/http"

	"coffee-site-server/internal/httpx"
	"coffee-site-server/internal/models"
)

func (h *Handler) ListUtensils(w http.ResponseWriter, r *http.Request) {
	utensils, err := h.Store.ListUtensils()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, utensils)
}

func (h *Handler) GetUtensil(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的商品編號")
		return
	}
	utensil, err := h.Store.GetUtensil(id)
	if err != nil {
		writeStoreError(w, err, "找不到商品")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, utensil)
}

type utensilRequest struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Image    string `json:"image"`
	Stock    int    `json:"stock"`
	Category string `json:"category"`
}

func (h *Handler) CreateUtensil(w http.ResponseWriter, r *http.Request) {
	var req utensilRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.Name == "" {
		httpx.Error(w, http.StatusBadRequest, "無效的商品資料")
		return
	}
	utensil, err := h.Store.CreateUtensil(models.Utensil{
		Name: req.Name, Price: req.Price, Image: req.Image, Stock: req.Stock, Category: req.Category,
	})
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "新增商品失敗")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, utensil)
}

func (h *Handler) UpdateUtensil(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的商品編號")
		return
	}
	var req utensilRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.Name == "" {
		httpx.Error(w, http.StatusBadRequest, "無效的商品資料")
		return
	}
	utensil, err := h.Store.UpdateUtensil(id, models.Utensil{
		Name: req.Name, Price: req.Price, Image: req.Image, Stock: req.Stock, Category: req.Category,
	})
	if err != nil {
		writeStoreError(w, err, "找不到商品")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, utensil)
}

func (h *Handler) DeleteUtensil(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的商品編號")
		return
	}
	if err := h.Store.DeleteUtensil(id); err != nil {
		writeStoreError(w, err, "找不到商品")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"message": "已刪除"})
}

func (h *Handler) ListUtensilCategories(w http.ResponseWriter, r *http.Request) {
	options, err := h.Store.ListOptions("utensil_categories")
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, options)
}

func (h *Handler) AddUtensilCategory(w http.ResponseWriter, r *http.Request) {
	h.addOption(w, r, "utensil_categories")
}

func (h *Handler) ToggleUtensilCategory(w http.ResponseWriter, r *http.Request) {
	h.toggleOption(w, r, "utensil_categories")
}
