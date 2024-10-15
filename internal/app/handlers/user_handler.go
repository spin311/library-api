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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserResponse
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if user.FirstName == "" || user.LastName == "" {
		helpers.WriteErrorResponse(w, errors.New("first_name and last_name parameters are required"), http.StatusBadRequest)
		return
	}

	httpErr := services.CreateUser(user)
	if !models.IsHttpErrorEmpty(httpErr) {
		helpers.WriteHttpErrorResponse(w, httpErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode("User created successfully")
}

func GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, httpErr := services.GetUsers()
	if !models.IsHttpErrorEmpty(httpErr) {
		helpers.WriteHttpErrorResponse(w, httpErr)
		return
	}
	if len(users) == 0 {
		users = []models.UserResponse{}
	}
	jsonErr := json.NewEncoder(w).Encode(users)
	if jsonErr != nil {
		helpers.WriteErrorResponse(w, jsonErr, http.StatusInternalServerError)
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["userId"])
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if id <= 0 {
		helpers.WriteErrorResponse(w, errors.New("invalid identifier"), http.StatusBadRequest)
		return
	}
	user, httpErr := services.GetUser(id)
	if !models.IsHttpErrorEmpty(httpErr) {
		helpers.WriteHttpErrorResponse(w, httpErr)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		helpers.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

}
