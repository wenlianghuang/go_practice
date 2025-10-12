package models

import "time"

type Book struct {
	ID        string     `json:"id"`
	Title     string     `json:"title" binding:"required"`
	Author    string     `json:"author" binding:"required"`
	Price     float64    `json:"price" binding:"required,min=0"`
	ISBN      string     `json:"isbn,omitempty"`
	Category  string     `json:"category,omitempty"`
	Published *time.Time `json:"published,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type BookPatch struct {
	Title     *string    `json:"title,omitempty"`
	Author    *string    `json:"author,omitempty"`
	Price     *float64   `json:"price,omitempty" binding:"omitempty,min=0"`
	ISBN      *string    `json:"isbn,omitempty"`
	Category  *string    `json:"category,omitempty"`
	Published *time.Time `json:"published,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
