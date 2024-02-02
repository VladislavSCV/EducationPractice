package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/VladislavSCV/EducationPractice"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	connString := "user=postgres password=31415926 dbname=SportMaster sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// Endpoints for users
	router.HandleFunc("/api/users", Controllers.GetUsers).Methods("GET") // Updated function name
	router.HandleFunc("/api/users/{id}", Controllers.GetUser).Methods("GET") // Updated function name

	// Endpoints for products
	router.HandleFunc("/api/products", Controllers.GetProducts).Methods("GET") // Updated function name
	router.HandleFunc("/api/products/{id}", Controllers.GetProduct).Methods("GET") // Updated function name

	// Endpoints for orders
	router.HandleFunc("/api/orders", Controllers.GetOrders).Methods("GET") // Updated function name
	router.HandleFunc("/api/orders/{id}", Controllers.GetOrder).Methods("GET") // Updated function name

	log.Fatal(http.ListenAndServe(":8000", router))
}
