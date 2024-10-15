package models

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Quantity      int    `json:"quantity"`
	BorrowedCount int    `json:"borrowed_count"`
}

// BookResponse represents a book in the system
//
//swagger:model
type BookResponse struct {
	//example: 1
	ID int `json:"id"`
	//example: The Great Gatsby
	Title string `json:"title"`
	//example: 5
	AvailableCount int `json:"quantity"`
}

func NewBookResponseFromBook(book Book) BookResponse {
	return BookResponse{
		ID:             book.ID,
		Title:          book.Title,
		AvailableCount: book.Quantity - book.BorrowedCount,
	}
}
