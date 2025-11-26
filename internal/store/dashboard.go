package store

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type DashboardStore struct {
	db *sql.DB
}

func (s *DashboardStore) Get(ctx context.Context, monthOffset int) (*models.Dashboard, error) {
	// Get the month range
	now := time.Now()
	start, end := utils.GetMonthRange(now, monthOffset)

	var wg sync.WaitGroup

	dashboard := &models.Dashboard{}
	errChan := make(chan error, 5)

	wg.Add(5)

	// Cement Summary
	go func() {
		defer wg.Done()
		summary, err := s.GetCementSummary(ctx, start, end)
		if err != nil {
			errChan <- err
			return
		}
		dashboard.CementSummary = summary
	}()

	// Sand Summary
	go func() {
		defer wg.Done()
		summary, err := s.GetSandSummary(ctx, start, end)
		if err != nil {
			errChan <- err
			return
		}
		dashboard.SandSummary = summary
	}()

	// Production Summary
	go func() {
		defer wg.Done()
		summary, err := s.GetProductionSummary(ctx, start, end)
		if err != nil {
			errChan <- err
			return
		}
		dashboard.ProductionSummary = summary
	}()

	// Transaction Summary
	go func() {
		defer wg.Done()
		summary, err := s.GetTransactionSummary(ctx, start, end)
		if err != nil {
			errChan <- err
			return
		}
		dashboard.TransactionSummary = summary
	}()

	// Salary Summary
	go func() {
		defer wg.Done()
		summary, err := s.GetSalarySummary(ctx, start, end)
		if err != nil {
			errChan <- err
			return
		}
		dashboard.SalarySummary = summary
	}()

	wg.Wait()
	close(errChan)

	// Check for errors
	if len(errChan) > 0 {
		return nil, <- errChan
	}

	return dashboard, nil
}

func (s *DashboardStore) GetCementSummary(ctx context.Context, start, end time.Time) (models.CementSummary, error) {
	var summary models.CementSummary
	query := `
		SELECT 
			COUNT(*) as total_stock,
			COALESCE(SUM(quantity), 0) as total_quantity,
			COALESCE(SUM(quantity * price_per_bag), 0) as total_price
		FROM cement_stocks
		WHERE purchase_date BETWEEN $1 AND $2
	`

	err := s.db.QueryRowContext(ctx, query, start, end).Scan(
		&summary.TotalStock,
		&summary.TotalQuantity,
		&summary.TotalPrice,
	)

	return summary, err
}

func (s *DashboardStore) GetSandSummary(ctx context.Context, start, end time.Time) (models.SandSummary, error) {
	var summary models.SandSummary
	query := `
		SELECT 
			COUNT(*) as total_purchase,
			COALESCE(SUM(quantity), 0) as total_quantity,
			COALESCE(SUM(quantity * price_per_truck), 0) as total_price
		FROM sand_purchases
		WHERE purchase_date BETWEEN $1 AND $2
	`

	err := s.db.QueryRowContext(ctx, query, start, end).Scan(
		&summary.TotalPurchase,
		&summary.TotalQuantity,
		&summary.TotalPrice,
	)

	return summary, err
}

func (s *DashboardStore) GetProductionSummary(ctx context.Context, start, end time.Time) (models.ProductionSummary, error) {
	var summary models.ProductionSummary
	query := `
		SELECT 
			COALESCE(SUM(quantity), 0) as total_production
		FROM productions
		WHERE production_date BETWEEN $1 AND $2
	`

	err := s.db.QueryRowContext(ctx, query, start, end).Scan(
		&summary.TotalProduction,
	)

	return summary, err
}

func (s *DashboardStore) GetTransactionSummary(ctx context.Context, start, end time.Time) (models.TransactionSummary, error) {
	var summary models.TransactionSummary
	query := `
		SELECT 
			COALESCE(SUM(quantity), 0) as total_transaction,
			COALESCE(SUM(total_price), 0) as total_income
		FROM transactions
		WHERE purchase_date BETWEEN $1 AND $2
	`

	err := s.db.QueryRowContext(ctx, query, start, end).Scan(
		&summary.TotalTransaction,
		&summary.TotalIncome,
	)

	return summary, err
}

func (s *DashboardStore) GetSalarySummary(ctx context.Context, start, end time.Time) (models.SalarySummary, error) {
	var summary models.SalarySummary
	query := `
		SELECT 
			COALESCE(SUM(salary), 0) as total_salary
		FROM employee_salary
		WHERE end_date BETWEEN $1 AND $2
	`

	err := s.db.QueryRowContext(ctx, query, start, end).Scan(
		&summary.TotalSalary,
	)

	return summary, err
}

