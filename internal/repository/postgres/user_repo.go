package postgres

import (
	"database/sql"
	"github.com/spin311/library-api/internal/repository/models"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func InsertUser(user models.User) error {
	_, err := db.Query(`INSERT INTO users (FIRST_NAME, LAST_NAME) VALUES ($1, $2)`, user.FirstName, user.LastName)
	if err != nil {
		return err
	}
	return nil
}

func GetUsers() ([]models.User, error) {
	rows, err := db.Query(`SELECT FIRST_NAME, LAST_NAME FROM users`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.FirstName, &user.LastName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
