package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	database "github.com/kevinbrivio/batako-backend/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env if available
	_ = godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer dbConn.Close()

	// Ensure schema is set
	if _, err := dbConn.Exec("SET search_path TO my_schema"); err != nil {
		log.Fatal("failed to set schema:", err)
	}

	if err := database.SeedAll(context.Background(), dbConn); err != nil {
		log.Fatalf("❌ seeding failed: %v", err)
	}

	log.Println("✅ seeding completed successfully")
}