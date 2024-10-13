package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spin311/library-api/internal/app/services"
	"net/http"
	"strconv"
)

func GetBooks(w http.ResponseWriter, _ *http.Request) {
	books, err := services.GetBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func BorrowBook(w http.ResponseWriter, r *http.Request) {
	userId, userErr := strconv.Atoi(r.URL.Query().Get("user_id"))
	bookId, bookErr := strconv.Atoi(r.URL.Query().Get("book_id"))
	if userErr != nil || bookErr != nil {
		http.Error(w, "Invalid user_id or book_id parameter", http.StatusBadRequest)
		return
	}
	err := services.BorrowBook(userId, bookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("book borrowed successfully")
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := services.GetBook(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	userId, userErr := strconv.Atoi(r.URL.Query().Get("user_id"))
	bookId, bookErr := strconv.Atoi(r.URL.Query().Get("book_id"))
	if userErr != nil || bookErr != nil {
		http.Error(w, "Invalid user_id or book_id parameter", http.StatusBadRequest)
		return
	}
	err := services.ReturnBook(userId, bookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("book returned successfully")
}
