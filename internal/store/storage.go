package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/kevinbrivio/batako-backend/internal/models"
)

type Storage struct {
	Production interface {
		Create(context.Context, *models.Production) error
		GetAll(context.Context, int, int) ([]models.Production, int, error)
		GetAllMonthly(context.Context, int) ([]models.Production, int, int, error)
		GetAllWeekly(context.Context, time.Time, time.Time) ([]models.Production, int, error)
		GetByID(context.Context, string) (*models.Production, error)
		GetTotalProduction(context.Context, time.Time, time.Time) (int, error)
		Update(context.Context, *models.Production) error
		Delete(context.Context, string) error
	}
	Transaction interface {
		Create(context.Context, *models.Transaction) error
		GetAll(context.Context, int, int) ([]models.Transaction, int, error)
		GetAllWeekly(context.Context, int) ([]models.Transaction, int, error)
		GetAllDaily(context.Context, time.Time) ([]models.Transaction, int, int, float64, error)
		GetAllMonthly(context.Context, int) ([]models.Transaction, int, int, float64, error)
		GetByID(context.Context, string) (*models.Transaction, error)
		Update(context.Context, *models.Transaction) error
		Delete(context.Context, string) error
		GetTotalWeeks(ctx context.Context) (int, error)
	}
	Salary interface {
		GetWeekly(context.Context, time.Time) (float64, int, error)
		GetMonthly(context.Context, int) ([]models.EmployeeSalary, int, error)
		AddSalary(context.Context, *models.EmployeeSalary) error
		GenerateWeeklySalary(context.Context) error
		StartSchedulers(context.Context)
	}
	CementStock interface {
		CheckTypes(context.Context, string) (int, error)
		Create(context.Context, *models.CementStock) error
		Update(context.Context, *models.CementStock) error
		GetAllMonthly(context.Context, int) ([]models.CementStock, int, int, float64, error)
		GetByType(context.Context, string, int) ([]models.CementStock, error)
		Delete(context.Context, string) error
	}
	SandPurchase interface {
		CheckTypes(context.Context, string) (int, error)
		Create(context.Context, *models.SandPurchase) error
		Update(context.Context, *models.SandPurchase) error
		GetAllMonthly(context.Context, int) ([]models.SandPurchase, int, int, float64, error)
		GetByType(context.Context, string, int) ([]models.SandPurchase, error)
		Delete(context.Context, string) error
	}
	Dashboard interface {
		Get(context.Context, int) (*models.Dashboard, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	prodStore := &ProductionStore{db: db}
	return Storage{
		Production:   &ProductionStore{db: db},
		Transaction:  &TransactionStore{db: db},
		Salary:       &SalaryStore{db: db, prodStore: prodStore},
		CementStock:  &CementStockStore{db: db},
		SandPurchase: &SandPurchaseStore{db: db},
		Dashboard: 		&DashboardStore{db: db},
	}
}
