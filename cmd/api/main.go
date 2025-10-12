package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	r.HandleFunc("/productions", prodHandler.CreateProduction).Methods("POST")
	
	log.Println("Server running at :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}