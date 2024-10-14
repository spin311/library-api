package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spin311/library-api/internal/app/handlers"
	"github.com/spin311/library-api/pkg/config"
	"log"
	"net/http"
)

func main() {

	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database: %v", err)
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

	r.HandleFunc("/users/{userId}/books/{bookId}/borrow", handlers.BorrowBook).Methods("POST")
	r.HandleFunc("/users/{userId}/books/{bookId}/return", handlers.ReturnBook).Methods("PUT")

	serverPort := config.GetEnvString("SERVER_PORT")
	log.Fatal(http.ListenAndServe(serverPort, r))
}
