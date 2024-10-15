package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spin311/library-api/internal/app/helpers"
	"github.com/spin311/library-api/internal/app/services"
	"github.com/spin311/library-api/internal/repository/models"
	"net/http"
	"strconv"
)

func GetBooks(w http.ResponseWriter, _ *http.Request) {
	books, err := services.GetBooks()
	if !models.IsHttpErrorEmpty(err) {
		helpers.WriteHttpErrorResponse(w, err)
		return
	}
	if len(books) == 0 {
		books = []models.BookResponse{}
	}
	jsonErr := json.NewEncoder(w).Encode(books)
	if jsonErr != nil {
		helpers.WriteErrorResponse(w, jsonErr, http.StatusInternalServerError)
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
	if !models.IsHttpErrorEmpty(err) {
		helpers.WriteHttpErrorResponse(w, err)
		return
	}
	jsonErr := json.NewEncoder(w).Encode(fmt.Sprintf("book with ID %d borrowed successfully", bookId))
	if jsonErr != nil {
		helpers.WriteErrorResponse(w, jsonErr, http.StatusInternalServerError)
		return
	}
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if id <= 0 {
		helpers.WriteHttpErrorResponse(w, models.NewHttpError("invalid identifier", http.StatusBadRequest))
		return
	}
	user, httpError := services.GetBook(id)
	if !models.IsHttpErrorEmpty(httpError) {
		helpers.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	jsonErr := json.NewEncoder(w).Encode(user)
	if jsonErr != nil {
		helpers.WriteErrorResponse(w, jsonErr, http.StatusInternalServerError)
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
	httpError := services.ReturnBook(userId, bookId)
	if !models.IsHttpErrorEmpty(httpError) {
		helpers.WriteHttpErrorResponse(w, httpError)
		return
	}
	jsonError := json.NewEncoder(w).Encode(fmt.Sprintf("book with ID %d returned successfully", bookId))
	if jsonError != nil {
		helpers.WriteErrorResponse(w, jsonError, http.StatusInternalServerError)
		return
	}
}
