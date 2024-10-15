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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with provided first name and last name
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {string} string "User created successfully"
// @Failure 400 {object} models.HttpError
// @Failure 500 {object} models.HttpError
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
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
	err = json.NewEncoder(w).Encode(fmt.Sprintf("User %s %s created successfully", user.FirstName, user.LastName))
}

// GetUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} models.HttpError
// @Router /users [get]
func GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, httpErr := services.GetUsers()
	if !models.IsHttpErrorEmpty(httpErr) {
		helpers.WriteHttpErrorResponse(w, httpErr)
		return
	}
	if len(users) == 0 {
		users = []models.User{}
	}
	jsonErr := json.NewEncoder(w).Encode(users)
	if jsonErr != nil {
		helpers.WriteErrorResponse(w, jsonErr, http.StatusInternalServerError)
		return
	}
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by user ID
// @Tags users
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.HttpError
// @Failure 404 {object} models.HttpError
// @Failure 500 {object} models.HttpError
// @Router /users/{userId} [get]
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
