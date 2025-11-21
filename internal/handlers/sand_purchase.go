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

type SandPurchaseHandler struct {
	Store store.Storage
}

func NewSandPurchaseHandler(s store.Storage) *SandPurchaseHandler{
	return &SandPurchaseHandler{Store: s}
}

func (h *SandPurchaseHandler) AddSandPurchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateSandPurchaseRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON format"))
		return
	}

	if req.Quantity <= 0 {
		utils.WriteError(w, utils.NewBadRequestError("Quantity is minimum 0."))
		return
	}

	domain := &models.SandPurchase{
		SandType: models.SandType{Name: req.SandTypeName},
		Quantity: req.Quantity,
		PricePerTruck: req.PricePerTruck,
		PurchaseDate: req.PurchaseDate,
	}
		
	if err := h.Store.SandPurchase.Create(ctx, domain); err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return 
	}

	utils.WriteJSON(w, http.StatusCreated, "Sand purchase successfully added", domain)
}

func (h *SandPurchaseHandler) UpdateSand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the ID from params
	idStr := chi.URLParam(r, "id")

	var stock models.UpdateSandPurchaseRequest
	if err := utils.ReadJSON(r, &stock); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON Format"))
		return
	}

	domain := &models.SandPurchase{
		ID: idStr,
		SandType: models.SandType{Name: stock.SandTypeName},
		Quantity: stock.Quantity,
		PricePerTruck: stock.PricePerTruck,
		PurchaseDate: stock.PurchaseDate,
	}

	err := h.Store.SandPurchase.Update(ctx, domain)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Sand purchases updated successfully", domain)
}

func (h *SandPurchaseHandler) GetMonthlySandPurchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get query params
	monthNum, _ := strconv.Atoi(r.URL.Query().Get("month"))
	currentMonth := int(time.Now().Month())
	targetOffset := monthNum - currentMonth

	if targetOffset < -6 {
		targetOffset += 12
	}

	stock, totalCount, totalQuantity, totalPrice, err := h.Store.SandPurchase.GetAllMonthly(ctx, targetOffset)
	
	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	data := map[string]interface{}{
		"total_count":  totalCount,
		"total_price":  totalPrice,
		"total_quantity":  totalQuantity,
    "month":        monthNum,
    "month_name":   time.Month(monthNum).String(),
    "stocks": stock,
    }

	utils.WriteJSON(w, http.StatusOK, "Sucessfully get monthly sand purchases", data)
}

func (h *SandPurchaseHandler) GetSandPurchaseByType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	typeName := chi.URLParam(r, "type")
	if typeName == "" {
		typeName = r.URL.Query().Get("type")
	}

	if typeName == "" {
		utils.WriteError(w, utils.NewBadRequestError("Type name is required"))
		return
	}

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

	stock, err := h.Store.SandPurchase.GetByType(ctx, typeName, targetMonth)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, utils.NewNotFoundError("Sand purchase not found"))
			return
		}

		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Successfully get sand purchase", stock)
}

func (h *SandPurchaseHandler) DeleteSandPurchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	err := h.Store.SandPurchase.Delete(ctx, idStr)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Sand purchase deleted successfully", nil)
}
