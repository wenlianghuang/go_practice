package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

var books = []Book{
	{ID: "1", Title: "1984", Author: "George Orwell", Price: 9.99},
	{ID: "2", Title: "Brave New World", Author: "Aldous Huxley", Price: 8.99},
}

func main() {
	router := gin.Default()

	router.GET("/books", getBooks)
	router.POST("/books", addBook)
	router.GET("/books/:id", getBookByID)

	router.Run("localhost:8080")
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func addBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookByID(c *gin.Context) {
	id := c.Param("id")

	for _, book := range books {
		if book.ID == id {
			c.IndentedJSON(http.StatusOK, book)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}
