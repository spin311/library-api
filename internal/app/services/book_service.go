package services

import (
	"errors"
	"github.com/spin311/library-api/internal/repository/models"
	"github.com/spin311/library-api/internal/repository/postgres"
)

func GetBooks() ([]models.BookResponse, error) {
	return postgres.GetBooks()
}

func BorrowBook(userId int, bookId int) error {
	if userId <= 0 || bookId <= 0 {
		return errors.New("invalid identifier values")
	}
	book, err := postgres.GetBook(bookId)
	if err != nil {
		return err
	}
	availableBooks := book.Quantity - book.BorrowedCount
	if availableBooks == 0 {
		return errors.New("no available books")
	}
	err = postgres.BorrowBook(userId, bookId, book.BorrowedCount+1)
	if err != nil {
		return err
	}
	return nil

}

func GetBook(id int) (models.BookResponse, error) {
	if id <= 0 {
		return models.BookResponse{}, errors.New("invalid identifier")
	}
	var bookResponse models.BookResponse
	book, err := postgres.GetBook(id)
	if err != nil {
		return bookResponse, err
	}
	bookResponse.Title = book.Title
	bookResponse.AvailableCount = book.Quantity - book.BorrowedCount
	return bookResponse, nil
}

func ReturnBook(userId int, bookId int) error {
	if userId <= 0 || bookId <= 0 {
		return errors.New("invalid identifier values")
	}
	book, err := postgres.GetBook(bookId)
	if err != nil {
		return err
	}
	if book.BorrowedCount == 0 {
		return errors.New("no borrowed books")
	}
	err = postgres.ReturnBook(userId, bookId, book.BorrowedCount-1)
	if err != nil {
		return err
	}
	return nil
}
