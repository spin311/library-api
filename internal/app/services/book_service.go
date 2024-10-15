package services

import (
	"github.com/spin311/library-api/internal/repository/models"
	"github.com/spin311/library-api/internal/repository/postgres"
	"net/http"
)

func GetBooks() ([]models.BookResponse, models.HttpError) {
	httpError := models.NewEmptyHttpError()
	books, httpError := postgres.GetBooks()
	return books, httpError
}

func BorrowBook(userId int, bookId int) models.HttpError {
	book, err := postgres.GetBook(bookId)
	if !models.IsHttpErrorEmpty(err) {
		return err
	}

	availableBooks := book.Quantity - book.BorrowedCount
	if availableBooks <= 0 {
		return models.NewHttpError("no available copies of the book", http.StatusBadRequest)
	}

	return postgres.BorrowBook(userId, bookId, book.BorrowedCount+1)
}

func GetBook(id int) (models.BookResponse, models.HttpError) {
	var bookResponse models.BookResponse
	book, err := postgres.GetBook(id)
	if !models.IsHttpErrorEmpty(err) {
		return bookResponse, err
	}
	bookResponse = models.NewBookResponseFromBook(book)
	return bookResponse, err
}

func ReturnBook(userId int, bookId int) models.HttpError {
	book, err := postgres.GetBook(bookId)
	if !models.IsHttpErrorEmpty(err) {
		return err
	}

	if book.BorrowedCount == 0 {
		return models.NewHttpError("no borrowed copies of this book to return", http.StatusBadRequest)
	}
	return postgres.ReturnBook(userId, bookId, book.BorrowedCount-1)
}
