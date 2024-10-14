package services

import (
	"errors"
	"fmt"
	"github.com/spin311/library-api/internal/repository/models"
	"github.com/spin311/library-api/internal/repository/postgres"
)

func GetBooks() ([]models.BookResponse, error) {
	return postgres.GetBooks()
}

func BorrowBook(userId int, bookId int) error {
	book, err := postgres.GetBook(bookId)
	if err != nil {
		return fmt.Errorf("failed to retrieve book: %w", err)
	}

	availableBooks := book.Quantity - book.BorrowedCount
	if availableBooks <= 0 {
		return errors.New("no available copies of the book")
	}

	err = postgres.BorrowBook(userId, bookId, book.BorrowedCount+1)
	if err != nil {
		return fmt.Errorf("failed to borrow book: %w", err)
	}

	return nil
}

func GetBook(id int) (models.BookResponse, error) {
	var bookResponse models.BookResponse
	book, err := postgres.GetBook(id)
	if err != nil {
		return bookResponse, fmt.Errorf("failed to retrieve book: %w", err)
	}

	bookResponse.Title = book.Title
	bookResponse.AvailableCount = book.Quantity - book.BorrowedCount
	return bookResponse, nil
}

func ReturnBook(userId int, bookId int) error {
	book, err := postgres.GetBook(bookId)
	if err != nil {
		return fmt.Errorf("failed to retrieve book: %w", err)
	}

	if book.BorrowedCount == 0 {
		return errors.New("no borrowed copies of this book to return")
	}

	err = postgres.ReturnBook(userId, bookId, book.BorrowedCount-1)
	if err != nil {
		return fmt.Errorf("failed to return book: %w", err)
	}

	return nil
}
