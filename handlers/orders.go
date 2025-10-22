package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shirazJafri/gomem-orders/models"
)

type OrderHandler struct {
	store *models.Store
}

func NewOrderHandler(store *models.Store) *OrderHandler {
	return &OrderHandler{store: store}
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orders := h.store.GetOrders()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Failed to encode orders", http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}

	order, exists := h.store.GetOrder(id)
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode order", http.StatusInternalServerError)
		return
	}
}
