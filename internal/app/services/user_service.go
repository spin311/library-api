package services

import (
	"github.com/spin311/library-api/internal/repository/models"
	"github.com/spin311/library-api/internal/repository/postgres"
)

func CreateUser(user models.UserResponse) models.HttpError {
	return postgres.InsertUser(user)
}

func GetUsers() ([]models.User, models.HttpError) {
	return postgres.GetUsers()
}

func GetUser(id int) (models.User, models.HttpError) {
	return postgres.GetUser(id)
}
