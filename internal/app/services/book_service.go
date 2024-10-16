package services

import (
	"fmt"
	"github.com/spin311/library-api/internal/repository/models"
	"github.com/spin311/library-api/internal/repository/postgres"
	"net/http"
)

func GetBooks() ([]models.BookResponse, models.HttpError) {
	return postgres.GetBooks()
}

func BorrowBook(userId int, bookId int) models.HttpError {
	return postgres.BorrowBook(userId, bookId)
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
		return models.NewHttpError(fmt.Sprintf("no borrowed copies exist for the book with ID %d", book.ID), http.StatusBadRequest)
	}
	return postgres.ReturnBook(userId, bookId, book.BorrowedCount-1)
}
