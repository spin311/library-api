package services

import (
	"github.com/spin311/library-api/internal/repository/models"
	"github.com/spin311/library-api/internal/repository/postgres"
)

func CreateUser(user models.User) error {
	return postgres.InsertUser(user)
}

func GetUsers() ([]models.User, error) {
	return postgres.GetUsers()
}
