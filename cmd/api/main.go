package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/kevinbrivio/batako-backend/internal/handlers"
	"github.com/kevinbrivio/batako-backend/internal/store"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// OPEN DB
	var connStr = os.Getenv("DATABASE_URL")
	log.Printf("Connecting to: %s (local dev or prod)", connStr[:len("postgres://")+10]) // Truncated log

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open DB: ", err.Error())
	}
	defer db.Close()

	// Ping to verify live connection
	err = db.Ping()
	if err != nil {
		log.Fatal("DB connection failed (ping): ", err.Error())
	}
	log.Println("Successfully connected and pinged DB")

	_, err = db.Exec("SET search_path TO my_schema")
	if err != nil {
		log.Fatal("Schema set failed: ", err.Error())
	}
	log.Println("Schema set to my_schema")

	storage := store.NewStorage(db)
	prodHandler := handlers.NewProductionHandler(storage)
	transactionHandler := handlers.NewTransactionHandler(storage)
	salaryHandler := handlers.NewSalaryStorage(storage)
	cementStockHandler := handlers.NewCementStockHandler(storage)
	sandPurchaseHandler := handlers.NewSandPurchaseHandler(storage)
	// Run scheduler -> Calculate weekly salary
	salaryStorage := storage.Salary
	go salaryStorage.StartSchedulers(context.Background())

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/productions", func(r chi.Router) {
		r.Post("/", prodHandler.CreateProduction)
		r.Get("/", prodHandler.GetAllProductions)
		r.Get("/monthly", prodHandler.GetProductionMonthly)
		r.Get("/weekly", prodHandler.GetProductionWeekly)
		r.Get("/{id}", prodHandler.GetProduction)
		r.Put("/{id}", prodHandler.UpdateProduction)
		r.Delete("/{id}", prodHandler.DeleteProduction)
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", transactionHandler.CreateTransaction)
		r.Get("/", transactionHandler.GetAllTransactions)
		r.Get("/daily", transactionHandler.GetTransactionsDaily)
		r.Get("/weekly", transactionHandler.GetTransactionsWeekly)
		r.Get("/monthly", transactionHandler.GetTransactionsMonthly)
		r.Get("/{id}", transactionHandler.GetTransaction)
		r.Put("/{id}", transactionHandler.UpdateTransaction)
		r.Delete("/{id}", transactionHandler.DeleteTransaction)
	})

	r.Route("/salary", func(r chi.Router) {
		r.Get("/weekly", salaryHandler.GetWeeklySalary)
		r.Get("/monthly", salaryHandler.GetMonthlySalaries)
	})

	r.Route("/stock", func(r chi.Router) {
		r.Route("/cement", func(r chi.Router) {
			r.Post("/", cementStockHandler.AddCementStock)
			r.Put("/{id}", cementStockHandler.UpdateCementStock)
			r.Get("/monthly", cementStockHandler.GetMonthlyCementStock)
			r.Get("/{type}", cementStockHandler.GetCementStocksByType)
			r.Delete("/{id}", cementStockHandler.DeleteCementStock)
		})
		r.Route("/sand", func(r chi.Router) {
			r.Post("/", sandPurchaseHandler.AddSandPurchase)
			r.Put("/{id}", sandPurchaseHandler.UpdateSand)
			// r.Get("/monthly", sandPurchaseHandler.GetMonthlySandPurchase)
			// r.Delete("{id}", sandPurchaseHandler.DeleteCementStock)
		})
	})

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
