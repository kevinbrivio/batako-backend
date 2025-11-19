package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type SandPurchaseStore struct {
	db *sql.DB
}

func (s *SandPurchaseStore) Create(ctx context.Context, sand *models.SandPurchase) error {
	sand.ID = uuid.New().String()

	query := `
		INSERT INTO sand_purchases (id, sand_type_id, quantity, price_per_truck, purchase_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	sandTypeID, err := s.CheckTypes(ctx, sand.SandType)
	if err != nil {
		return err
	}

	err = s.db.QueryRowContext(
		ctx,
		query,
		sand.ID,
		sandTypeID,
		sand.Quantity,
		sand.PricePerTruck,
		sand.PurchaseDate,
	).Scan(
		&sand.CreatedAt,
		&sand.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *SandPurchaseStore) Update(ctx context.Context, sand *models.SandPurchase) error {
	query := `
		UPDATE sand_purchases
		SET sand_type_id = $2, quantity = $3, price_per_truck = $4, purchase_date = $5
		WHERE id = $1
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	sandTypeID, err := s.CheckTypes(ctx, sand.SandType)
	if err != nil {
		return err
	}

	err = s.db.QueryRowContext(
		ctx,
		query,
		sand.ID,
		sandTypeID,
		sand.Quantity,
		sand.PricePerTruck,
		sand.PurchaseDate,
	).Scan(
		&sand.CreatedAt,
		&sand.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.NewNotFoundError("Sand purchase")
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *SandPurchaseStore) GetAllMonthly(ctx context.Context, monthOffset int) ([]models.SandPurchase, int, int, float64, error) {
	today := time.Now()

	start, end := utils.GetMonthRange(today, monthOffset)
	
	query := `
		SELECT 
			id, 
			sand_type_id, 
			quantity,
			price_per_truck,
			purchase_date,
			COUNT(*) OVER() as total_count,
			SUM(quantity) OVER() as total_quantity,
			SUM(price_per_truck) OVER() as total_price,
			created_at,
			updated_at
		FROM sand_purchases
		WHERE purchase_date BETWEEN $1 and $2
		ORDER BY purchase_date DESC
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	defer rows.Close()

	stocks := []models.SandPurchase{}
	var totalCount, totalQuantity int
	var totalPrice float64

	for rows.Next() {
		var stock models.SandPurchase
		if err := rows.Scan(
			&stock.ID,
			&stock.SandType,
			&stock.Quantity,
			&stock.PricePerTruck,
			&stock.PurchaseDate,
			&totalCount, 
			&totalQuantity,
			&totalPrice,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		); err != nil {
			return stocks, 0, 0, 0, err
		}
		stocks = append(stocks, stock)
	}
	if err = rows.Err(); err != nil {
		return stocks, 0, 0, 0, err
	}

	return stocks, totalCount, totalQuantity, totalPrice, nil
}

func (s *SandPurchaseStore) CheckTypes(ctx context.Context, typeName string) (int, error) {
	query := `
		SELECT id
		FROM sand_types
		WHERE name = $1
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	var id int
	err := s.db.QueryRowContext(ctx, query, typeName).Scan(&id)

	if err == sql.ErrNoRows {
		return -1, utils.NewNotFoundError("Sand type")
	}

	if err != nil {
		return -1, err
	}

	return id, nil
} 