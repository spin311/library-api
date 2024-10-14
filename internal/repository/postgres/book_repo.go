package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/spin311/library-api/internal/repository/models"
)

var dbBook *sql.DB

func SetBookDB(database *sql.DB) {
	dbBook = database
}

func GetBooks() ([]models.BookResponse, error) {
	stmt, err := dbBook.Prepare(`SELECT TITLE, QUANTITY, BORROWED_COUNT FROM books`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
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
		if err := rows.Scan(&book.Title, &quantity, &borrowedCount); err != nil {
			return nil, err
		}
		book.AvailableCount = quantity - borrowedCount
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func GetBook(bookId int) (models.Book, error) {
	var book models.Book
	stmt, err := dbBook.Prepare(`SELECT ID, TITLE, QUANTITY, BORROWED_COUNT FROM books WHERE id = $1`)
	if err != nil {
		return book, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt)

	if err := stmt.QueryRow(bookId).Scan(&book.ID, &book.Title, &book.Quantity, &book.BorrowedCount); err != nil {
		return book, err
	}

	return book, nil
}

func BorrowBook(userId int, bookId int, newCount int) error {
	ctx := context.Background()
	tx, err := dbBook.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmtBorrow, err := tx.PrepareContext(ctx, `INSERT INTO borrow (user_id, book_id) VALUES ($1, $2)`)
	if err != nil {
		_ = tx.Rollback()
		return err
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
		return err
	}

	err = updateBookCountWithTx(tx, bookId, newCount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
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
	return execErr
}

func ReturnBook(userId int, bookId int, newCount int) error {
	ctx := context.Background()
	tx, err := dbBook.BeginTx(ctx, nil)
	if err != nil {
		return err
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
		return err
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
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if rowsAffected == 0 {
		_ = tx.Rollback()
		return errors.New("no borrowed books found for this user")
	}

	err = updateBookCountWithTx(tx, bookId, newCount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
