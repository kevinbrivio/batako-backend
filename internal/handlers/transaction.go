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

type TransactionHandler struct {
	Store store.Storage
}

func NewTransactionHandler(s store.Storage) *TransactionHandler {
	return &TransactionHandler{Store: s}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.Transaction
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON format"))
		return
	}

	if req.Customer == "" {
		utils.WriteError(w, utils.NewBadRequestError("Customer name cannot be empty"))
		return 
	}

	if req.Address == "" {
		utils.WriteError(w, utils.NewBadRequestError("Address cannot be empty"))
		return 
	}

	if req.Quantity <= 0 {
		utils.WriteError(w, utils.NewBadRequestError("Quantity is minimum 0."))
		return
	}

	now := time.Now()
	if req.PurchaseDate.After(now) {
		utils.WriteError(w, utils.NewBadRequestError("Date cannot be in the future"))
		return
	}

	if err := h.Store.Transaction.Create(ctx, &req); err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Transaction created successfully", req)
}

func (h *TransactionHandler) GetTransactionsWeekly(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get query params
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // default to 1
	}

	weekOffset := page - 1

	t, totalCount, err := h.Store.Transaction.GetAllWeekly(ctx, weekOffset)
	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	totalPages, err := h.Store.Transaction.GetTotalWeeks(ctx)
	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
	}

	response := utils.PaginatedResponse{
		Items:      t,
		Total:      totalCount,
		Page:       page,
		PageSize:   len(t),
		TotalPages: totalPages,
	}

	utils.WriteJSON(w, http.StatusOK, "Sucessfully get weekly Transactions", response)
}

func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get query params
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1 // default to 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 6 // default to 6 -> weekly
	}

	// Calculate offset
	offset := (page - 1) * limit

	t, totalCount, err := h.Store.Transaction.GetAll(ctx, limit, offset)
	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	response := utils.PaginatedResponse{
		Items:      t,
		Total:      totalCount,
		Page:       page,
		PageSize:   limit,
		TotalPages: totalPages,
	}

	utils.WriteJSON(w, http.StatusOK, "Sucessfully get all Transactions", response)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the ID from params
	idStr := chi.URLParam(r, "id")
	// id, err := uuid.Parse(idStr)
	// if err != nil {
	// 	utils.WriteError(w, utils.NewBadRequestError("Invalid ID format"))
	// }

	t, err := h.Store.Transaction.GetByID(ctx, idStr)

	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Sucessfully get Transaction", t)
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the ID from params
	idStr := chi.URLParam(r, "id")

	var t models.Transaction
	if err := utils.ReadJSON(r, &t); err != nil {
		utils.WriteError(w, utils.NewBadRequestError("Invalid JSON Format"))
		return
	}

	t.ID = idStr

	err := h.Store.Transaction.Update(ctx, &t)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Transaction updated successfully", t)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	err := h.Store.Transaction.Delete(ctx, idStr)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Transaction deleted successfully", nil)
}
