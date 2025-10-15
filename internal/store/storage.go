package store

import (
	"context"
	"database/sql"
	"github.com/kevinbrivio/batako-backend/internal/models"
)

type Storage struct {
	Production interface {
		Create(context.Context, *models.Production) error
		GetAll(context.Context, int, int) ([]models.Production, error)
		GetByID(context.Context, string) (*models.Production, error)
		Update(context.Context, *models.Production) error
		Delete(context.Context, string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Production: &ProductionStore{db: db},
	}
}
