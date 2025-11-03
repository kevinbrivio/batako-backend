package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kevinbrivio/batako-backend/internal/store"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type SalaryHandler struct {
	Store store.Storage
}

func NewSalaryStorage(s store.Storage) *SalaryHandler {
	return &SalaryHandler{Store: s}
}

func (h *SalaryHandler) GetWeeklySalary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	dateStr := r.URL.Query().Get("date")
	var dt time.Time
	var err error
	if dateStr != "" {
		dt, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			dt = time.Now()
		}
	} else {
		dt = time.Now()
	}

	log.Println(dateStr)

	salary, totalProduction, err := h.Store.Salary.GetWeekly(ctx, dt)
	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	data := map[string]interface{}{
		"salary": salary,
		"total_production": totalProduction,
	}

	utils.WriteJSON(w, http.StatusOK, "Successfuly get weekly salary", data)
}

func (h *SalaryHandler) GetMonthlySalaries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get query params
	monthNum, _ := strconv.Atoi(r.URL.Query().Get("month"))
	currentMonth := int(time.Now().Month())
	targetOffset := monthNum - currentMonth

	if targetOffset < -6 {  
		targetOffset += 12  
	}

	salaries, totalCount, err := h.Store.Salary.GetMonthly(ctx, targetOffset)
	
	if err != nil {
		utils.WriteError(w, utils.NewInternalServerError(err))
		return
	}

	data := map[string]interface{}{
        "total_count":  totalCount,
        "month":        monthNum, 
        "month_name":   time.Month(monthNum).String(),
        "salaries": salaries,
    }


	utils.WriteJSON(w, http.StatusOK, "Sucessfully get monthly salaries", data)
}
