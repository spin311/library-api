package models

type Borrow struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	BookID     int    `json:"book_id"`
	BorrowedAt string `json:"borrowed_at"`
	ReturnedAt string `json:"returned_at"`
}
