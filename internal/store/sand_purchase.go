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

type SandPurchaseStore struct {
	db *sql.DB
}

func (s *SandPurchaseStore) Create(ctx context.Context, sand *models.SandPurchase) error {
	sand.ID = uuid.New().String()
	
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	sandTypeId, err := s.CheckTypes(ctx, sand.SandType.Name)
	if err != nil {
		return err
	}

	if sandTypeId == -1 {
		return utils.NewNotFoundError("Sand type with name: " + sand.SandType.Name)
	}

	query := `
		INSERT INTO sand_purchases (id, sand_type_id, quantity, price_per_truck, purchase_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING sand_type_id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	err = s.db.QueryRowContext(
		ctx,
		query,
		sand.ID,
		sandTypeId,
		sand.Quantity,
		sand.PricePerTruck,
		sand.PurchaseDate,
	).Scan(
		&sand.SandType.ID,
		&sand.CreatedAt,
		&sand.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return utils.NewConflictError("Sand purchase with this data already exists")
			}
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return utils.NewBadRequestError("Invalid sand type ID")
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SandPurchaseStore) Update(ctx context.Context, sand *models.SandPurchase) error {
	query := `
		UPDATE sand_purchases
		SET sand_type_id = $2, quantity = $3, price_per_truck = $4, purchase_date = $5
		WHERE id = $1
		RETURNING sand_type_id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	sandTypeID, err := s.CheckTypes(ctx, sand.SandType.Name)
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
		&sand.SandType.ID,
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
		st.id as sand_type_id,
		st.name as sand_type_name,
		SUM(sp.quantity) as total_quantity,
		SUM(sp.price_per_truck * sp.quantity) as total_price,
		AVG(sp.price_per_truck) as avg_price_per_truck,
		MIN(sp.purchase_date) as first_purchase_date,
		MAX(sp.purchase_date) as last_purchase_date
	FROM sand_purchases sp 
	JOIN sand_types st ON sp.sand_type_id = st.id
	WHERE sp.purchase_date BETWEEN $1 AND $2	
	GROUP BY st.id, st.name
	ORDER BY st.name ASC
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
		var totalPricePerType float64
		var totalQuantityPerType int

		if err = rows.Scan(
			&stock.SandType.ID,
			&stock.SandType.Name,
			&totalQuantityPerType,
			&totalPricePerType,
			&stock.PricePerTruck,
			&stock.PurchaseDate,
			&stock.UpdatedAt,
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

func (s *SandPurchaseStore) GetByType(ctx context.Context, typeName string, monthOffset int) ([]models.SandPurchase, error) {
	today := time.Now()
	start, end := utils.GetMonthRange(today, monthOffset)

	query := `
		SELECT 
			sp.id,
			st.id,
			st.name,
			sp.quantity,
			sp.price_per_truck,
			sp.purchase_date,
			sp.created_at,
			sp.updated_at
		FROM sand_purchases sp
		JOIN sand_types st 
		ON sp.sand_type_id = st.id 
		WHERE st.name = $1 AND sp.purchase_date BETWEEN $2 AND $3
		ORDER BY sp.purchase_date DESC
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, typeName, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.SandPurchase
	for rows.Next() {
		var stock models.SandPurchase
		err := rows.Scan(
			&stock.ID,
			&stock.SandType.ID,
			&stock.SandType.Name,
			&stock.Quantity,
			&stock.PricePerTruck,
			&stock.PurchaseDate,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		stock.TotalPrice = float64(stock.Quantity) * stock.PricePerTruck
		stocks = append(stocks, stock)
	}

	if len(stocks) == 0 {
		return []models.SandPurchase{}, nil
	}

	return stocks, nil
}

func (s *SandPurchaseStore) Delete(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		DELETE FROM sand_purchases
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return utils.NewNotFoundError("Sand purchase")
	}
	if err != nil {
		return err
	}

	return nil
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
