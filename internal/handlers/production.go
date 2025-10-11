package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/store"
)

type ProductionHandler struct {
	Store store.Storage
}

func NewProductionHandler(s store.Storage) *ProductionHandler{
	return &ProductionHandler{Store: s}
}

func (h *ProductionHandler) CreateProduction(w http.ResponseWriter, r *http.Request) {
	var req models.Production
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	if err := h.Store.Production.Create(ctx, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}