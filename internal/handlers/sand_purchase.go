package handlers

import (
	"net/http"
	"strconv"
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

	var req models.SandPurchase
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON format"))
		return
	}

	if req.Quantity <= 0 {
		utils.WriteError(w, utils.NewBadRequestError("Quantity is minimum 0."))
		return
	}

	if err := h.Store.SandPurchase.Create(ctx, &req); err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return 
	}

	utils.WriteJSON(w, http.StatusCreated, "Sand purchase successfully added",req)
}

func (h *SandPurchaseHandler) UpdateSand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the ID from params
	idStr := chi.URLParam(r, "id")

	var stock models.SandPurchase
	if err := utils.ReadJSON(r, &stock); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON Format"))
		return
	}

	stock.ID = idStr

	err := h.Store.SandPurchase.Update(ctx, &stock)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Sand purchases updated successfully", stock)
}

func (h *TransactionHandler) GetMonthlySandPurchase(w http.ResponseWriter, r *http.Request) {
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

