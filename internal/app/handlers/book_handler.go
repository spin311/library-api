package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/spin311/library-api/internal/app/helpers"
	"github.com/spin311/library-api/internal/app/services"
	"github.com/spin311/library-api/internal/repository/models"
	"net/http"
	"strconv"
)

func GetBooks(w http.ResponseWriter, _ *http.Request) {
	books, err := services.GetBooks()
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if len(books) == 0 {
		books = []models.BookResponse{}
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func BorrowBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, bookErr := strconv.Atoi(vars["bookId"])
	userId, userErr := strconv.Atoi(vars["userId"])
	if userErr != nil || bookErr != nil {
		helpers.WriteErrorResponse(w, errors.New("invalid userId or bookId parameter"), http.StatusBadRequest)
		return
	}
	if userId <= 0 || bookId <= 0 {
		helpers.WriteErrorResponse(w, errors.New("invalid identifier values"), http.StatusBadRequest)
		return
	}
	err := services.BorrowBook(userId, bookId)
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("book borrowed successfully")
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if id <= 0 {
		helpers.WriteErrorResponse(w, errors.New("invalid identifier"), http.StatusBadRequest)
		return
	}
	user, err := services.GetBook(id)
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, bookErr := strconv.Atoi(vars["bookId"])
	userId, userErr := strconv.Atoi(vars["userId"])
	if userErr != nil || bookErr != nil {
		helpers.WriteErrorResponse(w, errors.New("invalid userId or bookId parameter"), http.StatusBadRequest)
		return
	}
	if userId <= 0 || bookId <= 0 {
		helpers.WriteErrorResponse(w, errors.New("invalid identifier values"), http.StatusBadRequest)
		return
	}
	err := services.ReturnBook(userId, bookId)
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("book returned successfully")
}
