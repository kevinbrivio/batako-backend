package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type CementStockStore struct {
	db *sql.DB
}

func (s *CementStockStore) Create(ctx context.Context, c *models.CementStock) error {
	c.ID = uuid.New().String()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	cementTypeID, err := s.CheckTypes(ctx, c.CementType.Name)
	if err != nil {
		return err
	}

	if cementTypeID == -1 {
		return utils.NewNotFoundError("Cement type with name: " + c.CementType.Name)
	}

	query := `
		INSERT INTO cement_stocks (id, cement_type_id, quantity, price_per_bag, purchase_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err = tx.QueryRowContext(
		ctx,
		query,
		c.ID,
		cementTypeID,
		c.Quantity,
		c.PricePerBag,
		c.PurchaseDate,
	).Scan(
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return utils.NewConflictError("Cement stock with this data already exists")
			}
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return utils.NewBadRequestError("Invalid cement type ID")
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *CementStockStore) Update(ctx context.Context, c *models.CementStock) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE cement_stocks
		SET cement_type_id = $2, quantity = $3, price_per_bag = $4, purchase_date = $5
		WHERE id = $1
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	cementTypeID, err := s.CheckTypes(ctx, c.CementType.Name)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(
		ctx,
		query,
		c.ID,
		cementTypeID,
		c.Quantity,
		c.PricePerBag,
		c.PurchaseDate,
	).Scan(
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return utils.NewConflictError("Cement stock with this data already exists")
			}
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return utils.NewBadRequestError("Invalid cement type ID")
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *CementStockStore) GetAllMonthly(ctx context.Context, monthOffset int) ([]models.CementStock, int, int, float64, error) {
	today := time.Now()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	start, end := utils.GetMonthRange(today, monthOffset)

	query := `
	SELECT 
		ct.id AS cement_type_id,
		ct.name AS cement_type_name,
		SUM(cs.quantity) AS total_quantity,
		SUM(cs.price_per_bag * cs.quantity) AS total_price,
		AVG(cs.price_per_bag) AS avg_price_per_bag,
		MIN(cs.purchase_date) AS first_purchase_date,
		MAX(cs.purchase_date) AS last_purchase_date
		FROM cement_stocks cs
		JOIN cement_types ct ON cs.cement_type_id = ct.id
		WHERE cs.purchase_date BETWEEN $1 AND $2
		GROUP BY ct.id, ct.name
		ORDER BY ct.name ASC;
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	defer rows.Close()

	stocks := []models.CementStock{}
	var totalCount, totalQuantity int
	var totalPrice float64

	for rows.Next() {
		var stock models.CementStock
		var totalPricePerType float64
		var totalQuantityPerType int
		if err := rows.Scan(
			&stock.CementType.ID,
			&stock.CementType.Name,
			&totalQuantityPerType, // total_quantity
			&totalPricePerType,    // total_price
			&stock.PricePerBag,    // avg_price_per_bag (optional)
			&stock.PurchaseDate,   // first_purchase_date (or drop if not needed)
			&stock.UpdatedAt,      // last_purchase_date (optional)
		); err != nil {
			return stocks, 0, 0, 0, err
		}
		stock.Quantity = totalQuantityPerType
		stock.TotalPrice = totalPricePerType
		stocks = append(stocks, stock)

		totalQuantity += totalQuantityPerType
		totalPrice += totalPricePerType
		totalCount++
	}
	if err = rows.Err(); err != nil {
		return stocks, 0, 0, 0, err
	}

	return stocks, totalCount, totalQuantity, totalPrice, nil
}

func (s *CementStockStore) CheckTypes(ctx context.Context, typeName string) (int, error) {
	query := `
		SELECT id
		FROM cement_types
		WHERE name = $1
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var id int
	err := s.db.QueryRowContext(ctx, query, typeName).Scan(&id)

	if err == sql.ErrNoRows {
		return -1, utils.NewNotFoundError("Cement type")
	}

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *CementStockStore) GetByType(ctx context.Context, typeName string, monthOffset int) ([]models.CementStock, error) {
	today := time.Now()
	start, end := utils.GetMonthRange(today, monthOffset)

	query := `
		SELECT 
			cs.id,
			ct.id,
			ct.name,
			cs.quantity,
			cs.price_per_bag,
			cs.purchase_date,
			cs.created_at,
			cs.updated_at
		FROM cement_stocks cs
		JOIN cement_types ct
		ON cs.cement_type_id = ct.id
		WHERE ct.name = $1 AND cs.purchase_date BETWEEN $2 AND $3
		ORDER BY cs.purchase_date DESC
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, typeName, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.CementStock

	for rows.Next() {
		var stock models.CementStock
		err := rows.Scan(
			&stock.ID,
			&stock.CementType.ID,
			&stock.CementType.Name,
			&stock.Quantity,
			&stock.PricePerBag,
			&stock.PurchaseDate,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		stock.TotalPrice = float64(stock.Quantity) * stock.PricePerBag
		stocks = append(stocks, stock)
	}

	if len(stocks) == 0 {
		return []models.CementStock{}, nil
	}

	return stocks, nil
}

func (s *CementStockStore) Delete(ctx context.Context, stockId string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		DELETE FROM cement_stocks
		WHERE id = $1;	
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, stockId)
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
