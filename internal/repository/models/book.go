package models

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Quantity      int    `json:"quantity"`
	BorrowedCount int    `json:"borrowed_count"`
}

type BookResponse struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	AvailableCount int    `json:"quantity"`
}

func NewBookResponseFromBook(book Book) BookResponse {
	return BookResponse{
		ID:             book.ID,
		Title:          book.Title,
		AvailableCount: book.Quantity - book.BorrowedCount,
	}
}
