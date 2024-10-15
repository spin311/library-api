package models

// User represents a user in the system
//
//swagger:model
type User struct {
	//example: 1
	ID int `json:"id"`
	//example: John
	FirstName string `json:"first_name"`
	//example: Doe
	LastName string `json:"last_name"`
}

type UserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
