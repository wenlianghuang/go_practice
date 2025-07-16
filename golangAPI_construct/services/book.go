package services

import (
	"errors"
	"golangAPI_construct/models"
	"strconv"
)

var books = []models.Book{
	{ID: "1", Title: "1984", Author: "George Orwell", Price: 9.99},
	{ID: "2", Title: "Brave New World", Author: "Aldous Huxley", Price: 8.99},
	{ID: "3", Title: "To Kill a Mockingbird", Author: "Harper Lee", Price: 12.99},
}

var nextID = 4

type BookService struct{}

func NewBookService() *BookService {
	return &BookService{}
}

func (s *BookService) GetAllBooks() []models.Book {
	return books
}

func (s *BookService) GetBooksByAuthor(author string) []models.Book {
	var filteredBooks []models.Book
	for _, book := range books {
		if book.Author == author {
			filteredBooks = append(filteredBooks, book)
		}
	}
	return filteredBooks
}

func (s *BookService) GetBookByID(id string) (*models.Book, error) {
	for _, book := range books {
		if book.ID == id {
			return &book, nil
		}
	}
	return nil, errors.New("book not found")
}

func (s *BookService) CreateBook(book models.Book) (*models.Book, error) {
	// 檢查是否已存在相同 ID
	for _, existingBook := range books {
		if existingBook.ID == book.ID {
			return nil, errors.New("book with this ID already exists")
		}
	}

	// 如果沒有提供 ID，自動生成
	if book.ID == "" {
		book.ID = strconv.Itoa(nextID)
		nextID++
	}

	books = append(books, book)
	return &book, nil
}

func (s *BookService) UpdateBook(id string, book models.Book) (*models.Book, error) {
	for i, existingBook := range books {
		if existingBook.ID == id {
			book.ID = id
			books[i] = book
			return &book, nil
		}
	}
	return nil, errors.New("book not found")
}

func (s *BookService) PatchBook(id string, patch models.BookPatch) (*models.Book, error) {
	for i, book := range books {
		if book.ID == id {
			if patch.Title != nil {
				books[i].Title = *patch.Title
			}
			if patch.Author != nil {
				books[i].Author = *patch.Author
			}
			if patch.Price != nil {
				books[i].Price = *patch.Price
			}
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func (s *BookService) DeleteBook(id string) (*models.Book, error) {
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			return &book, nil
		}
	}
	return nil, errors.New("book not found")
}

func (s *BookService) GetBooksCount() int {
	return len(books)
}
