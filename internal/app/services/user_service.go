package services

import (
	"errors"
	"github.com/spin311/library-api/internal/repository/models"
	"github.com/spin311/library-api/internal/repository/postgres"
)

func CreateUser(user models.UserResponse) error {
	if user.FirstName == "" || user.LastName == "" {
		return errors.New("first_name and last_name parameters are required")
	}

	return postgres.InsertUser(user)
}

func GetUsers() ([]models.UserResponse, error) {
	return postgres.GetUsers()
}

func GetUser(id int) (models.UserResponse, error) {
	if id <= 0 {
		return models.UserResponse{}, errors.New("invalid identifier")
	}

	return postgres.GetUser(id)
}
