package services

import (
	"fmt"
	"github.com/spin311/library-api/internal/repository/models"
	"github.com/spin311/library-api/internal/repository/postgres"
)

func CreateUser(user models.UserResponse) error {
	err := postgres.InsertUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func GetUsers() ([]models.UserResponse, error) {
	users, err := postgres.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

func GetUser(id int) (models.UserResponse, error) {
	user, err := postgres.GetUser(id)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return user, nil
}
