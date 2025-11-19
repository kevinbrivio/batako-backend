package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/store"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type CementStockHandler struct {
	Store store.Storage
}

func NewCementStockHandler(s store.Storage) *CementStockHandler {
	return &CementStockHandler{Store: s}
}

func (h *CementStockHandler) AddCementStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateCementStockRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON format"))
		return
	}

	if req.Quantity <= 0 {
		utils.WriteError(w, utils.NewBadRequestError("Quantity is minimum 0."))
		return
	}

	if req.PricePerBag <= 0 {
		utils.WriteError(w, utils.NewBadRequestError("Price per bag must be greater than 0"))
		return
	}

	if req.PurchaseDate.After(time.Now()) {
		utils.WriteError(w, utils.NewBadRequestError("Purchase date cannot be in the future."))
		return
	}

	domain := &models.CementStock{
		CementType:   models.CementType{Name: req.CementTypeName},
		Quantity:     req.Quantity,
		PricePerBag:  req.PricePerBag,
		PurchaseDate: req.PurchaseDate,
	}

	if err := h.Store.CementStock.Create(ctx, domain); err != nil {
		if utils.IsNotFound(err) {
			utils.WriteError(w, utils.NewNotFoundError("Cement type not found"))
			return
		}
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Cement stock successfully added", req)
}

func (h *CementStockHandler) UpdateCementStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the ID from params
	idStr := chi.URLParam(r, "id")

	var stock models.UpdateCementStockRequest
	if err := utils.ReadJSON(r, &stock); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON Format"))
		return
	}

	if stock.Quantity <= 0 {
		utils.WriteError(w, utils.NewBadRequestError("Quantity is minimum 0."))
		return
	}

	if stock.PurchaseDate.After(time.Now()) {
		utils.WriteError(w, utils.NewBadRequestError("Purchase date cannot be in the future."))
		return
	}

	domain := &models.CementStock{
		ID:           idStr,
		CementType:   models.CementType{Name: stock.CementTypeName},
		Quantity:     stock.Quantity,
		PricePerBag:  stock.PricePerBag,
		PurchaseDate: stock.PurchaseDate,
	}

	err := h.Store.CementStock.Update(ctx, domain)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Cement stock updated successfully", stock)
}

func (h *CementStockHandler) GetMonthlyCementStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get query params
	monthNum, _ := strconv.Atoi(r.URL.Query().Get("month"))
	currentMonth := int(time.Now().Month())
	targetOffset := monthNum - currentMonth

	if targetOffset < -6 {
		targetOffset += 12
	}

	stock, totalCount, totalQuantity, totalPrice, err := h.Store.CementStock.GetAllMonthly(ctx, targetOffset)

	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	data := map[string]interface{}{
		"total_count":    totalCount,
		"total_price":    totalPrice,
		"total_quantity": totalQuantity,
		"month":          monthNum,
		"month_name":     time.Month(monthNum).String(),
		"stocks":         stock,
	}

	utils.WriteJSON(w, http.StatusOK, "Sucessfully get monthly cement stock", data)
}

func (h *CementStockHandler) GetCementStocksByType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get type from path or query
	typeName := chi.URLParam(r, "type")
	if typeName == "" {
		typeName = r.URL.Query().Get("type")
	}

	// Validate type name
	if typeName == "" {
		utils.WriteError(w, utils.NewBadRequestError("Type name is required"))
		return
	}

	// Get month from query params with validation
	monthStr := r.URL.Query().Get("month")
	var targetMonth int
	if monthStr != "" {
		targetMonth, err := strconv.Atoi(monthStr)
		if err != nil || targetMonth < 1 || targetMonth > 12 {
			utils.WriteError(w, utils.NewBadRequestError("Invalid month. Must be between 1-12"))
			return
		}
	} else {
		targetMonth = int(time.Now().Month())
	}

	stock, err := h.Store.CementStock.GetByType(ctx, typeName, targetMonth)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, utils.NewNotFoundError("Cement stock not found"))
			return
		}

		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Sucessfully get cement stock", stock)
}

func (h *CementStockHandler) DeleteCementStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	err := h.Store.CementStock.Delete(ctx, idStr)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Transaction deleted successfully", nil)
}
