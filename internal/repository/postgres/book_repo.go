package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/spin311/library-api/internal/repository/models"
	"net/http"
)

var dbBook *sql.DB

func SetBookDB(database *sql.DB) {
	dbBook = database
}

func GetBooks() ([]models.BookResponse, models.HttpError) {
	stmt, err := dbBook.Prepare(`SELECT ID, TITLE, QUANTITY, BORROWED_COUNT FROM books`)
	if err != nil {
		return nil, models.NewHttpErrorFromError("failed to prepare statement", err, http.StatusInternalServerError)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	rows, err := stmt.Query()
	if err != nil {
		return nil, models.NewHttpErrorFromError("failed to execute query", err, http.StatusInternalServerError)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var books []models.BookResponse
	for rows.Next() {
		var book models.BookResponse
		var quantity, borrowedCount int
		if err := rows.Scan(&book.ID, &book.Title, &quantity, &borrowedCount); err != nil {
			return nil, models.NewHttpErrorFromError("failed to scan row", err, http.StatusInternalServerError)
		}
		book.AvailableCount = quantity - borrowedCount
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, models.NewHttpErrorFromError("rows error", err, http.StatusInternalServerError)
	}

	return books, models.NewEmptyHttpError()
}

func GetBook(bookId int) (models.Book, models.HttpError) {
	var book models.Book
	stmt, err := dbBook.Prepare(`SELECT ID, TITLE, QUANTITY, BORROWED_COUNT FROM books WHERE id = $1`)
	if err != nil {
		return book, models.NewHttpErrorFromError("failed to prepare statement", err, http.StatusInternalServerError)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	if err := stmt.QueryRow(bookId).Scan(&book.ID, &book.Title, &book.Quantity, &book.BorrowedCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return book, models.NewHttpError(fmt.Sprintf("book with ID %d not found", bookId), http.StatusNotFound)
		}
		return book, models.NewHttpErrorFromError("failed to scan row", err, http.StatusInternalServerError)
	}

	return book, models.NewEmptyHttpError()
}

func BorrowBook(userId int, bookId int) models.HttpError {
	ctx := context.Background()
	tx, err := dbBook.BeginTx(ctx, nil)
	if err != nil {
		return models.NewHttpErrorFromError("failed to begin transaction", err, http.StatusInternalServerError)
	}

	// Lock the row for the book to prevent race conditions
	stmtLock, err := tx.PrepareContext(ctx, `SELECT quantity, borrowed_count FROM books WHERE id = $1 FOR UPDATE`)
	if err != nil {
		_ = tx.Rollback()
		return models.NewHttpErrorFromError("failed to prepare lock statement", err, http.StatusInternalServerError)
	}
	defer func(stmtLock *sql.Stmt) {
		err := stmtLock.Close()
		if err != nil {
			return
		}
	}(stmtLock)

	var quantity, borrowedCount int
	err = stmtLock.QueryRow(bookId).Scan(&quantity, &borrowedCount)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return models.NewHttpError(fmt.Sprintf("book with ID %d not found", bookId), http.StatusNotFound)
		}
		return models.NewHttpErrorFromError("failed to scan book row", err, http.StatusInternalServerError)
	}

	availableBooks := quantity - borrowedCount
	if availableBooks <= 0 {
		_ = tx.Rollback()
		return models.NewHttpError(fmt.Sprintf("no available copies of the book with ID %d", bookId), http.StatusConflict)
	}

	stmtBorrow, err := tx.PrepareContext(ctx, `INSERT INTO borrow (user_id, book_id) VALUES ($1, $2)`)
	if err != nil {
		_ = tx.Rollback()
		return models.NewHttpErrorFromError("failed to prepare borrow statement", err, http.StatusInternalServerError)
	}
	defer func(stmtBorrow *sql.Stmt) {
		err := stmtBorrow.Close()
		if err != nil {
			return
		}
	}(stmtBorrow)

	_, err = stmtBorrow.Exec(userId, bookId)
	if err != nil {
		_ = tx.Rollback()
		return models.NewHttpErrorFromError("failed to execute borrow statement", err, http.StatusInternalServerError)
	}

	updateErr := updateBookCountWithTx(tx, bookId, borrowedCount+1)
	if updateErr != nil {
		_ = tx.Rollback()
		return models.NewHttpErrorFromError("failed to update book count", updateErr, http.StatusInternalServerError)
	}

	if err := tx.Commit(); err != nil {
		return models.NewHttpErrorFromError("failed to commit transaction", err, http.StatusInternalServerError)
	}

	return models.NewEmptyHttpError()
}

func updateBookCountWithTx(tx *sql.Tx, bookId int, newCount int) error {
	stmtUpdate, err := tx.PrepareContext(context.Background(), `UPDATE books SET borrowed_count = $1 WHERE id = $2`)
	if err != nil {
		return err
	}
	defer func(stmtUpdate *sql.Stmt) {
		err := stmtUpdate.Close()
		if err != nil {
			return
		}
	}(stmtUpdate)

	_, execErr := stmtUpdate.Exec(newCount, bookId)
	if execErr != nil {
		return execErr
	}
	return nil
}

func ReturnBook(userId int, bookId int, newCount int) models.HttpError {
	ctx := context.Background()
	tx, err := dbBook.BeginTx(ctx, nil)
	if err != nil {
		return models.NewHttpErrorFromError("failed to begin transaction", err, http.StatusInternalServerError)
	}

	stmtReturn, err := tx.PrepareContext(ctx, `
		WITH borrowed AS (
			SELECT id 
			  FROM borrow
			WHERE book_id = $1 
			    AND user_id = $2 
			    AND returned_at IS NULL
			ORDER BY borrowed_at
			LIMIT 1
		)
		UPDATE borrow
			   SET returned_at = CURRENT_TIMESTAMP
		 WHERE id IN (SELECT id FROM borrowed)
	`)
	if err != nil {
		_ = tx.Rollback()
		return models.NewHttpErrorFromError("failed to prepare statement", err, http.StatusInternalServerError)
	}
	defer func(stmtReturn *sql.Stmt) {
		err := stmtReturn.Close()
		if err != nil {
			return
		}
	}(stmtReturn)

	result, err := stmtReturn.Exec(bookId, userId)
	if err != nil {
		_ = tx.Rollback()
		return models.NewHttpErrorFromError("failed to execute statement", err, http.StatusInternalServerError)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return models.NewHttpErrorFromError("failed to get rows affected", err, http.StatusInternalServerError)
	}
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return models.NewHttpError(fmt.Sprintf("no borrowed books found for user ID %d and book ID %d", userId, bookId), http.StatusBadRequest)
	}

	err = updateBookCountWithTx(tx, bookId, newCount)
	if err != nil {
		_ = tx.Rollback()
		return models.NewHttpErrorFromError("failed to update book count", err, http.StatusInternalServerError)
	}

	if err := tx.Commit(); err != nil {
		return models.NewHttpErrorFromError("failed to commit transaction", err, http.StatusInternalServerError)
	}

	return models.NewEmptyHttpError()
}
