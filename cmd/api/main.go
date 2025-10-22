package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/kevinbrivio/batako-backend/internal/handlers"
	"github.com/kevinbrivio/batako-backend/internal/store"
	_ "github.com/lib/pq"
)

func main() {
	// OPEN DB
	var connStr = os.Getenv("CONN_STR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("SET search_path TO my_schema")
    if err != nil {
        log.Fatal(err)
    }

	storage := store.NewStorage(db)
	prodHandler := handlers.NewProductionHandler(storage)
	transactionHandler := handlers.NewTransactionHandler(storage)

	r := chi.NewRouter()
	r.Route("/productions", func(r chi.Router) {
		r.Post("/", prodHandler.CreateProduction)
		r.Get("/", prodHandler.GetAllProductions)
		r.Get("/{id}", prodHandler.GetProduction)
		r.Put("/{id}", prodHandler.UpdateProduction)
		r.Delete("/{id}", prodHandler.DeleteProduction)
	})
	
	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", transactionHandler.CreateTransaction)
		r.Get("/", transactionHandler.GetAllTransactions)
		r.Get("/weekly", transactionHandler.GetTransactionsWeekly)
		r.Get("/monthly", transactionHandler.GetTransactionsMonthly)
		r.Get("/{id}", transactionHandler.GetTransaction)
		r.Put("/{id}", transactionHandler.UpdateTransaction)
		r.Delete("/{id}", transactionHandler.DeleteTransaction)
	})
	
	log.Println("Server running at :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
