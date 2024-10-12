package postgres

import (
	"context"
	"database/sql"
	"github.com/spin311/library-api/internal/repository/models"
	"log"
)

var dbBook *sql.DB

func SetBookDB(database *sql.DB) {
	dbBook = database
}

func GetBooks() ([]models.BookResponse, error) {
	rows, err := dbBook.Query(`SELECT TITLE, QUANTITY, BORROWED_COUNT FROM books`)
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
	count := dbBook.QueryRow(`SELECT ID, TITLE, QUANTITY, BORROWED_COUNT FROM books WHERE id = $1`, bookId)
	if err := count.Scan(&book.ID, &book.Title, &book.Quantity, &book.BorrowedCount); err != nil {
		return book, err
	}
	return book, nil
}

func BorrowBook(userId int, bookId int, newCount int) error {
	ctx := context.Background()
	tx, err := dbBook.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = updateBookCountWithTx(tx, bookId, newCount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO borrow (user_id, book_id) VALUES ($1, $2)`, userId, bookId)
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
	_, execErr := tx.ExecContext(context.Background(), `UPDATE books SET borrowed_count = $1 WHERE id = $2`, newCount, bookId)
	return execErr
}
