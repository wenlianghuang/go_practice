// struct to json
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Book struct {
	Title  string `json:"title"`
	Author Author `json:"author"`
}

type Author struct {
	Sales     int  `json:"book_sales"`
	Age       int  `json:"age"`
	Developer bool `json:"is_developer"`
}

func main() {
	author := Author{
		Sales:     3,
		Age:       25,
		Developer: true,
	}

	book := Book{Title: "Learning Go", Author: author}

	// Convert the struct to JSON, and have it indented
	byteSlice, _ := json.MarshalIndent(book, "", "  ") // Add indentation here
	fmt.Println(string(byteSlice))
	err := os.WriteFile("book.json", byteSlice, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
	}
}
