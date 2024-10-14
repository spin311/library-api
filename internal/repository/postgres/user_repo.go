package postgres

import (
	"database/sql"
	"errors"
	"github.com/spin311/library-api/internal/repository/models"
)

var dbUser *sql.DB

func SetUserDB(database *sql.DB) {
	dbUser = database
}

func InsertUser(user models.UserResponse) error {
	stmt, err := dbUser.Prepare(`INSERT INTO users (FIRST_NAME, LAST_NAME) VALUES ($1, $2)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.FirstName, user.LastName)
	return err
}

func GetUsers() ([]models.UserResponse, error) {
	rows, err := dbUser.Query(`SELECT FIRST_NAME, LAST_NAME FROM users`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
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

func GetUser(id int) (models.UserResponse, error) {
	var user models.UserResponse
	stmt, err := dbUser.Prepare(`SELECT FIRST_NAME, LAST_NAME FROM users WHERE ID = $1`)
	if err != nil {
		return user, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	row := stmt.QueryRow(id)
	if err := row.Scan(&user.FirstName, &user.LastName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}
