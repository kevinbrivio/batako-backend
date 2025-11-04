package database

import (
	"context"
	"database/sql"
	"fmt"
)

// SeedAll runs all seeding routines safely and idempotently.
func SeedAll(ctx context.Context, db *sql.DB) error {
	if err := seedCementTypes(ctx, db); err != nil {
		return fmt.Errorf("cement_types: %w", err)
	}
	if err := seedSandTypes(ctx, db); err != nil {
		return fmt.Errorf("sand_types: %w", err)
	}
	return nil
}

func seedCementTypes(ctx context.Context, db *sql.DB) error {
	cementTypes := []string{"Tiga Roda", "Conch", "Merdeka", "Padang", "Rajawali"}

	for _, name := range cementTypes {
		_, err := db.ExecContext(ctx, `
			INSERT INTO my_schema.cement_types (name)
			VALUES ($1)
			ON CONFLICT (name) DO NOTHING;
		`, name)
		if err != nil {
			return fmt.Errorf("insert %s: %w", name, err)
		}
	}
	return nil
}

func seedSandTypes(ctx context.Context, db *sql.DB) error {
	sandTypes := []string{"Putih", "Kuning"}

	for _, name := range sandTypes {
		_, err := db.ExecContext(ctx, `
			INSERT INTO my_schema.sand_types (name)
			VALUES ($1)
			ON CONFLICT (name) DO NOTHING;
		`, name)
		if err != nil {
			return fmt.Errorf("insert %s: %w", name, err)
		}
	}
	return nil
}