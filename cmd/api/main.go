package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "github.com/spin311/library-api/docs"
	"github.com/spin311/library-api/internal/app/handlers"
	"github.com/spin311/library-api/pkg/config"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Library API
// @version 1.0
// @description Simple library API
// @host localhost:8080
// @BasePath /
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
	r.HandleFunc("/users", handlers.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users", handlers.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{userId}", handlers.GetUser).Methods(http.MethodGet)

	//Book Routes
	r.HandleFunc("/books", handlers.GetBooks).Methods(http.MethodGet)
	r.HandleFunc("/books/{bookId}", handlers.GetBook).Methods(http.MethodGet)

	r.HandleFunc("/users/{userId}/books/{bookId}/borrow", handlers.BorrowBook).Methods(http.MethodPost)
	r.HandleFunc("/users/{userId}/books/{bookId}/return", handlers.ReturnBook).Methods(http.MethodPut)

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	serverPort := config.GetEnvString("SERVER_PORT")
	log.Fatal(http.ListenAndServe(serverPort, r))
}
