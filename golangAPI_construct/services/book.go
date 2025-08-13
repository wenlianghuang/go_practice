package services

import (
	"errors"
	"golangAPI_construct/models"
	"strconv"
	"sync"
)

var (
	// 統一錯誤
	ErrBookNotFound = errors.New("book not found")
)

type BookService struct {
	mu     sync.RWMutex
	books  []models.Book
	nextID int
}

// 初始化內建資料
func NewBookService() *BookService {
	return &BookService{
		books: []models.Book{
			{ID: "1", Title: "1984", Author: "George Orwell", Price: 9.99},
			{ID: "2", Title: "Brave New World", Author: "Aldous Huxley", Price: 8.99},
			{ID: "3", Title: "To Kill a Mockingbird", Author: "Harper Lee", Price: 12.99},
		},
		nextID: 4,
	}
}

func (s *BookService) GetAllBooks() []models.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := make([]models.Book, len(s.books))
	copy(cp, s.books)
	return cp
}

func (s *BookService) GetBooksByAuthor(author string) []models.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var filtered []models.Book
	for _, b := range s.books {
		if b.Author == author {
			filtered = append(filtered, b)
		}
	}
	return filtered
}

func (s *BookService) GetBookByID(id string) (*models.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.books {
		if s.books[i].ID == id {
			return &s.books[i], nil
		}
	}
	return nil, ErrBookNotFound
}

func (s *BookService) CreateBook(book models.Book) (*models.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 忽略外部傳入的 ID，統一由系統產生
	book.ID = strconv.Itoa(s.nextID)
	s.nextID++

	s.books = append(s.books, book)
	return &s.books[len(s.books)-1], nil
}

func (s *BookService) UpdateBook(id string, book models.Book) (*models.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.books {
		if s.books[i].ID == id {
			book.ID = id
			s.books[i] = book
			return &s.books[i], nil
		}
	}
	return nil, ErrBookNotFound
}

func (s *BookService) PatchBook(id string, patch models.BookPatch) (*models.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.books {
		if s.books[i].ID == id {
			if patch.Title != nil {
				s.books[i].Title = *patch.Title
			}
			if patch.Author != nil {
				s.books[i].Author = *patch.Author
			}
			if patch.Price != nil {
				s.books[i].Price = *patch.Price
			}
			return &s.books[i], nil
		}
	}
	return nil, ErrBookNotFound
}

func (s *BookService) DeleteBook(id string) (*models.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.books {
		if s.books[i].ID == id {
			deleted := s.books[i]
			s.books = append(s.books[:i], s.books[i+1:]...)
			return &deleted, nil
		}
	}
	return nil, ErrBookNotFound
}

func (s *BookService) GetBooksCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.books)
}
