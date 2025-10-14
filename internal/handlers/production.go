package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/store"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type ProductionHandler struct {
	Store store.Storage
}

func NewProductionHandler(s store.Storage) *ProductionHandler{
	return &ProductionHandler{Store: s}
}

func (h *ProductionHandler) CreateProduction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.Production
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON format"))
		return
	}

	if req.Quantity <= 0 {
		utils.WriteError(w, utils.NewBadRequestError("Quantity is minimum 0."))
		return
	}
	
	
	if err := h.Store.Production.Create(ctx, &req); err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return 
	}

	utils.WriteJSON(w, http.StatusCreated, req)
}

func (h *ProductionHandler) GetAllProductions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	prods, err := h.Store.Production.GetAll(ctx);
	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}
	
	utils.WriteJSON(w, http.StatusOK, prods)
}

func (h *ProductionHandler) GetProduction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the ID from params
	idStr := chi.URLParam(r, "id")
	// id, err := uuid.Parse(idStr)
	// if err != nil {
	// 	utils.WriteError(w, utils.NewBadRequestError("Invalid ID format"))
	// }

	prod, err := h.Store.Production.GetByID(ctx, idStr)

	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, prod)
}

func (h *ProductionHandler) UpdateProduction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the ID from params
	idStr := chi.URLParam(r, "id")

	var prod models.Production
	if err := utils.ReadJSON(r, &prod); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON Format"))
		return
	}

	prod.ID = idStr

	err := h.Store.Production.Update(ctx, &prod)
	if err != nil {
		utils.WriteError(w, err)
		return 
	}

	utils.WriteJSON(w, http.StatusOK, prod)
}