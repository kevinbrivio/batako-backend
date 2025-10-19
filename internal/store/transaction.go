package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type TransactionStore struct {
	db *sql.DB
}

func (s *TransactionStore) Create(ctx context.Context, t *models.Transaction) error {
	t.ID = uuid.New().String()

	query := `
		INSERT INTO transactions (id, customer, address, quantity, total_price, purchase_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	// Calculate total price
	const price = 1600
	totalPrice := float64(t.Quantity) * price

	err := s.db.QueryRowContext(
		ctx,
		query,
		t.ID,
		t.Customer,
		t.Address,
		t.Quantity,
		totalPrice,
		t.PurchaseDate,
	).Scan(
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionStore) GetAll(ctx context.Context, limit, offset int) ([]models.Transaction, int, error) {
	query := `
		SELECT 
			id, 
			customer, 
			address,
			quantity,
			total_price,
			COUNT(*) OVER() as total_count,
			purchase_date,
			created_at,
			updated_at
		FROM transactions
		ORDER BY purchase_date DESC
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

	transactions := []models.Transaction{}
	var totalCount int

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.Customer,
			&t.Address,
			&t.Quantity,
			&t.TotalPrice,
			&totalCount, 
			&t.PurchaseDate,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return transactions, 0, err
		}
		transactions = append(transactions, t)
	}
	if err = rows.Err(); err != nil {
		return transactions, 0, err
	}

	return transactions, totalCount, nil
}

func (s *TransactionStore) GetByID(ctx context.Context, pID string) (*models.Transaction, error) {
	query := `
		SELECT * FROM transactions
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	var t models.Transaction

	err := s.db.QueryRowContext(
		ctx, query,
		pID,
	).Scan(
		&t.ID,
		&t.Customer,
		&t.Address,
		&t.Quantity,
		&t.TotalPrice,
		&t.PurchaseDate,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError("Transaction")
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *TransactionStore) Update(ctx context.Context, t *models.Transaction) error {
	query := `
		UPDATE transactions
		SET customer = $2, address = $3, quantity = $4, total_price = $4, purchase_date = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	// Calculate total price
	const price = 1600
	totalPrice := float64(t.Quantity) * price

	err := s.db.QueryRowContext(
		ctx,
		query,
		t.ID,
		t.Customer,
		t.Address,
		t.Quantity,
		totalPrice,
		t.PurchaseDate,
	).Scan(
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.NewNotFoundError("Transaction")
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionStore) Delete(ctx context.Context, tID string) error {
	query := `
		DELETE FROM transactions
		WHERE id = $1;	
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, tID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return utils.NewNotFoundError("Transaction")
	}
	return nil
}
