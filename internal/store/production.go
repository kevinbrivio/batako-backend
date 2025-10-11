package store

import (
	"context"
	"database/sql"
	"time"
	"github.com/kevinbrivio/batako-backend/internal/models"
)

type ProductionStore struct {
	db *sql.DB
}

func (s *ProductionStore) Create(ctx context.Context, p *models.Production) error {
	query := `
		INSERT INTO Productions (quantity, cement_used, sand_used, total)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		p.Quantity,
		p.CementUsed,
		p.SandUsed,
		p.Total,
	).Scan(
		&p.ID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
