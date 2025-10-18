package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type ProductionStore struct {
	db *sql.DB
}

func (s *ProductionStore) Create(ctx context.Context, p *models.Production) error {
	// Generate UUID before inserting
	p.ID = uuid.New().String()
	
	query := `
		INSERT INTO Productions (id, quantity, cement_used, date)
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
		p.Date,
	).Scan(
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *ProductionStore) GetAll(ctx context.Context, limit, offset int) ([]models.Production, int, error) {
	query := `
		SELECT 
			id, 
			quantity,
			cement_used,
			date,
			COUNT(*) OVER() as total_count,
			created_at,
			updated_at
		FROM productions
		ORDER BY date DESC
		LIMIT $1 OFFSET $2
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	// Pass limit and offset
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	productions := []models.Production{}
	var totalCount int

	for rows.Next() {
		var p models.Production
		if err := rows.Scan(
			&p.ID,
			&p.Quantity,
			&p.CementUsed,
			&p.Date,
			&totalCount, 
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return productions, 0, err
		}
		productions = append(productions, p)
	}
	if err = rows.Err(); err != nil {
		return productions, 0, err
	}

	return productions, totalCount, nil
}

func (s *ProductionStore) GetByID(ctx context.Context, pID string) (*models.Production, error) {
	query := `
		SELECT * FROM productions
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	var p models.Production

	err := s.db.QueryRowContext(
		ctx, query,
		pID,
	).Scan(
		&p.ID,
		&p.Quantity,
		&p.CementUsed,
		&p.Date,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError("Production")
	}
	
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *ProductionStore) Update(ctx context.Context, p *models.Production) error {
	query := `
		UPDATE productions
		SET quantity = $2, cement_used = $3, date = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		p.ID,
		p.Quantity,
		p.CementUsed,
		p.Date,
	).Scan(
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.NewNotFoundError("Production")
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *ProductionStore) Delete(ctx context.Context, pID string) error {
	query := `
		DELETE FROM productions
		WHERE id = $1;	
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, pID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return utils.NewNotFoundError("Production")
	}
	return nil
}