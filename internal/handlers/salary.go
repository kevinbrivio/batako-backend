package handlers

import (
	"net/http"

	"github.com/kevinbrivio/batako-backend/internal/store"
)

type SalaryHandler struct {
	Store store.Storage
}

func NewSalaryStorage(s store.Storage) *SalaryHandler {
	return &SalaryHandler{Store: s}
}

func (h *SalaryHandler) GetSalary(w http.ResponseWriter, r *http.Request) {
	
}
