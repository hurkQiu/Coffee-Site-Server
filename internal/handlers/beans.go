package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"coffee-site-server/internal/httpx"
	"coffee-site-server/internal/models"
	"coffee-site-server/internal/store"
)

func (h *Handler) ListBeans(w http.ResponseWriter, r *http.Request) {
	beans, err := h.Store.ListBeans()
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, beans)
}

func (h *Handler) GetBean(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的商品編號")
		return
	}
	bean, err := h.Store.GetBean(id)
	if err != nil {
		writeStoreError(w, err, "找不到商品")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, bean)
}

type beanRequest struct {
	Name    string `json:"name"`
	Price   int    `json:"price"`
	Image   string `json:"image"`
	Stock   int    `json:"stock"`
	Roast   string `json:"roast"`
	Process string `json:"process"`
}

func (h *Handler) CreateBean(w http.ResponseWriter, r *http.Request) {
	var req beanRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.Name == "" {
		httpx.Error(w, http.StatusBadRequest, "無效的商品資料")
		return
	}
	bean, err := h.Store.CreateBean(models.CoffeeBean{
		Name: req.Name, Price: req.Price, Image: req.Image, Stock: req.Stock, Roast: req.Roast, Process: req.Process,
	})
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "新增商品失敗")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, bean)
}

func (h *Handler) UpdateBean(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的商品編號")
		return
	}
	var req beanRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.Name == "" {
		httpx.Error(w, http.StatusBadRequest, "無效的商品資料")
		return
	}
	bean, err := h.Store.UpdateBean(id, models.CoffeeBean{
		Name: req.Name, Price: req.Price, Image: req.Image, Stock: req.Stock, Roast: req.Roast, Process: req.Process,
	})
	if err != nil {
		writeStoreError(w, err, "找不到商品")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, bean)
}

func (h *Handler) DeleteBean(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的商品編號")
		return
	}
	if err := h.Store.DeleteBean(id); err != nil {
		writeStoreError(w, err, "找不到商品")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"message": "已刪除"})
}

func (h *Handler) ListRoastLevels(w http.ResponseWriter, r *http.Request) {
	options, err := h.Store.ListOptions("roast_levels")
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, options)
}

func (h *Handler) AddRoastLevel(w http.ResponseWriter, r *http.Request) {
	h.addOption(w, r, "roast_levels")
}

func (h *Handler) ToggleRoastLevel(w http.ResponseWriter, r *http.Request) {
	h.toggleOption(w, r, "roast_levels")
}

func (h *Handler) ListProcessMethods(w http.ResponseWriter, r *http.Request) {
	options, err := h.Store.ListOptions("process_methods")
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, options)
}

func (h *Handler) AddProcessMethod(w http.ResponseWriter, r *http.Request) {
	h.addOption(w, r, "process_methods")
}

func (h *Handler) ToggleProcessMethod(w http.ResponseWriter, r *http.Request) {
	h.toggleOption(w, r, "process_methods")
}

type optionRequest struct {
	Name string `json:"name"`
}

func (h *Handler) addOption(w http.ResponseWriter, r *http.Request, table string) {
	var req optionRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.Name == "" {
		httpx.Error(w, http.StatusBadRequest, "無效的類別名稱")
		return
	}
	if err := h.Store.AddOption(table, req.Name); err != nil {
		if errors.Is(err, store.ErrConflict) {
			httpx.Error(w, http.StatusConflict, "此類別已存在")
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "新增類別失敗")
		return
	}
	options, err := h.Store.ListOptions(table)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, options)
}

func (h *Handler) toggleOption(w http.ResponseWriter, r *http.Request, table string) {
	name := chi.URLParam(r, "name")
	if err := h.Store.ToggleOptionHidden(table, name); err != nil {
		writeStoreError(w, err, "找不到類別")
		return
	}
	options, err := h.Store.ListOptions(table)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, options)
}

func parseIDParam(r *http.Request, key string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, key), 10, 64)
}

func writeStoreError(w http.ResponseWriter, err error, notFoundMessage string) {
	if errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusNotFound, notFoundMessage)
		return
	}
	httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
}
