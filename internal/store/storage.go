package store

import (
	"context"
	"database/sql"
	"github.com/kevinbrivio/batako-backend/internal/models"
)

type Storage struct {
	Production interface {
		Create(context.Context, *models.Production) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Production: &ProductionStore{db: db},
	}
}
