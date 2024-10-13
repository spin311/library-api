package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spin311/library-api/internal/app/handlers"
	"github.com/spin311/library-api/pkg/config"
	"net/http"

	"log"
)

const serverAddress = ":8080"

func main() {

	db := config.InitDatabase()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	config.SetDbs(db)
	r := mux.NewRouter()

	//User Routes
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{userId}", handlers.GetUser).Methods("GET")

	//Book Routes
	r.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	r.HandleFunc("/books/{bookId}", handlers.GetBook).Methods("GET")

	r.HandleFunc("/books/borrow", handlers.BorrowBook).Methods("POST")
	r.HandleFunc("/books/return", handlers.ReturnBook).Methods("POST")

	log.Fatal(http.ListenAndServe(serverAddress, r))
}
