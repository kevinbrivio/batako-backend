package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kevinbrivio/batako-backend/internal/handlers"
	"github.com/kevinbrivio/batako-backend/internal/store"
)

func main() {
	// OPEN DB
	var DB_ADDR = "postgres://batako_user@localhost:5432/batako?sslmode=disable"
	db, err := sql.Open("postgres", DB_ADDR)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	storage := store.NewStorage(db)
	prodHandler := handlers.NewProductionHandler(storage)

	r := mux.NewRouter()
	r.HandleFunc("/productions", prodHandler.CreateProduction).Methods("POST")
	
	log.Println("Server running at :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}