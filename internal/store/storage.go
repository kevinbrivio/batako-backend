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
		GetByID(context.Context, string) (*models.Production, error)
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
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Production: &ProductionStore{db: db},
		Transaction: &TransactionStore{db: db},
	}
}
