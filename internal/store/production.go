package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kevinbrivio/batako-backend/internal/models"
)

type ProductionStore struct {
	db *sql.DB
}

func (s *ProductionStore) Create(ctx context.Context, p *models.Production) error {
	// Generate UUID before inserting
	p.ID = uuid.New().String()
	
	query := `
		INSERT INTO Productions (id, quantity, cement_used, sand_used)
		VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		p.ID,
		p.Quantity,
		p.CementUsed,
		p.SandUsed,
	).Scan(
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
