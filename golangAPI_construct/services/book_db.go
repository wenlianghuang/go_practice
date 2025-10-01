package services

import (
	"context"
	"database/sql"
	"errors"
	"golangAPI_construct/models"
	"strconv"
)

type BookServiceDB struct {
	db *sql.DB
}

func NewBookServiceDB(db *sql.DB) BookServiceInterface {
	return &BookServiceDB{db: db}
}

// 編譯期保證 BookServiceDB 有實作介面
var _ BookServiceInterface = (*BookServiceDB)(nil)

func (s *BookServiceDB) GetAllBooks() []models.Book {
	rows, err := s.db.QueryContext(context.Background(),
		`SELECT id, title, author, price FROM books ORDER BY id`)
	if err != nil {
		return []models.Book{}
	}
	defer rows.Close()

	var list []models.Book
	for rows.Next() {
		var (
			id     int64
			title  string
			author string
			price  float64
		)
		if err := rows.Scan(&id, &title, &author, &price); err != nil {
			continue
		}
		list = append(list, models.Book{
			ID:     strconv.FormatInt(id, 10),
			Title:  title,
			Author: author,
			Price:  price,
		})
	}
	return list
}

func (s *BookServiceDB) GetBooksByAuthor(author string) []models.Book {
	rows, err := s.db.QueryContext(context.Background(),
		`SELECT id, title, author, price FROM books WHERE author = $1 ORDER BY id`, author)
	if err != nil {
		return []models.Book{}
	}
	defer rows.Close()

	var list []models.Book
	for rows.Next() {
		var (
			id    int64
			title string
			a     string
			price float64
		)
		if err := rows.Scan(&id, &title, &a, &price); err != nil {
			continue
		}
		list = append(list, models.Book{
			ID:     strconv.FormatInt(id, 10),
			Title:  title,
			Author: a,
			Price:  price,
		})
	}
	return list
}

func (s *BookServiceDB) GetBookByID(id string) (*models.Book, error) {
	nid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrBookNotFound
	}
	var (
		title  string
		author string
		price  float64
	)
	err = s.db.QueryRowContext(context.Background(),
		`SELECT title, author, price FROM books WHERE id = $1`, nid).
		Scan(&title, &author, &price)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBookNotFound
	}
	if err != nil {
		return nil, err
	}
	return &models.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Price:  price,
	}, nil
}

func (s *BookServiceDB) CreateBook(book models.Book) (*models.Book, error) {
	var lastID int64
	err := s.db.QueryRowContext(context.Background(),
		`INSERT INTO books (title, author, price) VALUES ($1,$2,$3) RETURNING id`,
		book.Title, book.Author, book.Price).Scan(&lastID)
	if err != nil {
		return nil, err
	}
	book.ID = strconv.FormatInt(lastID, 10)
	return &book, nil
}

func (s *BookServiceDB) UpdateBook(id string, book models.Book) (*models.Book, error) {
	nid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrBookNotFound
	}
	res, err := s.db.ExecContext(context.Background(),
		`UPDATE books SET title = $1, author = $2, price = $3 WHERE id = $4`,
		book.Title, book.Author, book.Price, nid)
	if err != nil {
		return nil, err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return nil, ErrBookNotFound
	}
	book.ID = id
	return &book, nil
}

func (s *BookServiceDB) PatchBook(id string, patch models.BookPatch) (*models.Book, error) {
	// 首先獲取現有記錄
	current, err := s.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	// 套用部分更新
	if patch.Title != nil {
		current.Title = *patch.Title
	}
	if patch.Author != nil {
		current.Author = *patch.Author
	}
	if patch.Price != nil {
		current.Price = *patch.Price
	}

	// 執行完整更新
	return s.UpdateBook(id, *current)
}

func (s *BookServiceDB) DeleteBook(id string) (*models.Book, error) {
	cur, err := s.GetBookByID(id)
	if err != nil {
		return nil, err
	}
	nid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrBookNotFound
	}
	res, err := s.db.ExecContext(context.Background(),
		`DELETE FROM books WHERE id = $1`, nid)
	if err != nil {
		return nil, err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return nil, ErrBookNotFound
	}
	return cur, nil
}

func (s *BookServiceDB) GetBooksCount() int {
	var n int
	if err := s.db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM books`).Scan(&n); err != nil {
		return 0
	}
	return n
}
