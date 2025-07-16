package models

type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title" binding:"required"`
	Author string  `json:"author" binding:"required"`
	Price  float64 `json:"price" binding:"required,min=0"`
}

type BookPatch struct {
	Title  *string  `json:"title,omitempty"`
	Author *string  `json:"author,omitempty"`
	Price  *float64 `json:"price,omitempty" binding:"omitempty,min=0"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
