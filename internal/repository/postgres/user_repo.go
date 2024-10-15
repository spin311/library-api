package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/spin311/library-api/internal/repository/models"
	"net/http"
)

var dbUser *sql.DB

func SetUserDB(database *sql.DB) {
	dbUser = database
}

func InsertUser(user models.UserResponse) models.HttpError {
	stmt, err := dbUser.Prepare(`INSERT INTO users (FIRST_NAME, LAST_NAME) VALUES ($1, $2)`)
	if err != nil {
		return models.NewHttpErrorFromError("failed to prepare statement", err, http.StatusInternalServerError)
	}
	_, err = stmt.Exec(user.FirstName, user.LastName)
	if err != nil {
		return models.NewHttpErrorFromError("failed to execute statement", err, http.StatusInternalServerError)
	}
	return models.NewEmptyHttpError()
}

func GetUsers() ([]models.User, models.HttpError) {
	rows, err := dbUser.Query(`SELECT ID, FIRST_NAME, LAST_NAME FROM users`)
	if err != nil {
		return nil, models.NewHttpErrorFromError("failed to query users", err, http.StatusInternalServerError)
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
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName); err != nil {
			return nil, models.NewHttpErrorFromError("failed to scan user", err, http.StatusInternalServerError)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, models.NewHttpErrorFromError("failed to iterate over users", err, http.StatusInternalServerError)
	}
	return users, models.NewEmptyHttpError()
}

func GetUser(id int) (models.User, models.HttpError) {
	var user models.User
	stmt, err := dbUser.Prepare(`SELECT ID, FIRST_NAME, LAST_NAME FROM users WHERE ID = $1`)
	if err != nil {
		return user, models.NewHttpErrorFromError("failed to prepare statement", err, http.StatusInternalServerError)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	row := stmt.QueryRow(id)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, models.NewHttpError(fmt.Sprintf("user with ID %d not found", id), http.StatusNotFound)
		}
		return user, models.NewHttpErrorFromError("failed to scan user", err, http.StatusInternalServerError)
	}
	return user, models.NewEmptyHttpError()
}
