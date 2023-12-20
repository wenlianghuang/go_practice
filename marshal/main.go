package main

import (
	"encoding/json"
	"fmt"
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

	byteSlice, _ := json.MarshalIndent(book, "", "")
	fmt.Println(string(byteSlice))
}
