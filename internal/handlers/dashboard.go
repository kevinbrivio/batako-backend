package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kevinbrivio/batako-backend/internal/store"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type DashboardHandler struct {
	Store store.Storage
}

func NewDashboardHandler(s store.Storage) *DashboardHandler {
	return &DashboardHandler{Store: s}
}

func (h *DashboardHandler) GetMonthly(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	currMonth := int(time.Now().Month())
	targetMonth := month - currMonth

	if targetMonth < -6 {
		targetMonth += 12
	}

	data, err := h.Store.Dashboard.Get(ctx, targetMonth)
	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Successfully get monthly dashboard", data)
}
