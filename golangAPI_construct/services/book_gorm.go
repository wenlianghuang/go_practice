package services

import (
	"errors"
	"fmt"
	"golangAPI_construct/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BookServiceGORM GORM 版本的書籍服務
// 提供更強大的數據庫操作功能，包括關聯查詢、統計和搜索
type BookServiceGORM struct {
	db *gorm.DB
}

// NewBookServiceGORM 創建 GORM 書籍服務實例
func NewBookServiceGORM(db *gorm.DB) BookServiceInterface {
	return &BookServiceGORM{db: db}
}

// 編譯期保證 BookServiceGORM 有實作介面
var _ BookServiceInterface = (*BookServiceGORM)(nil)

// GetAllBooks 獲取所有書籍
func (s *BookServiceGORM) GetAllBooks() []models.Book {
	var booksGORM []models.BookGORM
	if err := s.db.Find(&booksGORM).Error; err != nil {
		return []models.Book{}
	}

	// 轉換為原始模型
	books := make([]models.Book, len(booksGORM))
	for i, bookGORM := range booksGORM {
		books[i] = bookGORM.ToBook()
	}
	return books
}

// GetBooksByAuthor 根據作者獲取書籍
func (s *BookServiceGORM) GetBooksByAuthor(author string) []models.Book {
	var booksGORM []models.BookGORM
	if err := s.db.Where("author ILIKE ?", "%"+author+"%").Find(&booksGORM).Error; err != nil {
		return []models.Book{}
	}

	books := make([]models.Book, len(booksGORM))
	for i, bookGORM := range booksGORM {
		books[i] = bookGORM.ToBook()
	}
	return books
}

// GetBookByID 根據 ID 獲取書籍
func (s *BookServiceGORM) GetBookByID(id string) (*models.Book, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, ErrBookNotFound
	}

	var bookGORM models.BookGORM
	if err := s.db.First(&bookGORM, uint(idUint)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	book := bookGORM.ToBook()
	return &book, nil
}

// CreateBook 創建新書籍
func (s *BookServiceGORM) CreateBook(book models.Book) (*models.Book, error) {
	var bookGORM models.BookGORM
	bookGORM.FromBook(book)

	if err := s.db.Create(&bookGORM).Error; err != nil {
		return nil, err
	}

	createdBook := bookGORM.ToBook()
	return &createdBook, nil
}

// UpdateBook 更新書籍
func (s *BookServiceGORM) UpdateBook(id string, book models.Book) (*models.Book, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, ErrBookNotFound
	}

	var bookGORM models.BookGORM
	if err := s.db.First(&bookGORM, uint(idUint)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	// 更新欄位
	bookGORM.Title = book.Title
	bookGORM.Author = book.Author
	bookGORM.Price = book.Price

	if err := s.db.Save(&bookGORM).Error; err != nil {
		return nil, err
	}

	updatedBook := bookGORM.ToBook()
	return &updatedBook, nil
}

// PatchBook 部分更新書籍
func (s *BookServiceGORM) PatchBook(id string, patch models.BookPatch) (*models.Book, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, ErrBookNotFound
	}

	var bookGORM models.BookGORM
	if err := s.db.First(&bookGORM, uint(idUint)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	// 應用部分更新
	if patch.Title != nil {
		bookGORM.Title = *patch.Title
	}
	if patch.Author != nil {
		bookGORM.Author = *patch.Author
	}
	if patch.Price != nil {
		bookGORM.Price = *patch.Price
	}

	if err := s.db.Save(&bookGORM).Error; err != nil {
		return nil, err
	}

	updatedBook := bookGORM.ToBook()
	return &updatedBook, nil
}

// DeleteBook 刪除書籍（軟刪除）
func (s *BookServiceGORM) DeleteBook(id string) (*models.Book, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, ErrBookNotFound
	}

	var bookGORM models.BookGORM
	if err := s.db.First(&bookGORM, uint(idUint)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	book := bookGORM.ToBook()

	// 軟刪除
	if err := s.db.Delete(&bookGORM).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

// GetBooksCount 獲取書籍總數
func (s *BookServiceGORM) GetBooksCount() int {
	var count int64
	s.db.Model(&models.BookGORM{}).Count(&count)
	return int(count)
}

// ===================
// GORM 專有的高級功能
// ===================

// SearchBooks 搜索書籍（支持標題、作者、ISBN）
func (s *BookServiceGORM) SearchBooks(query string) ([]models.BookSearchResult, error) {
	var results []models.BookSearchResult

	searchPattern := "%" + strings.ToLower(query) + "%"

	err := s.db.Model(&models.BookGORM{}).
		Select("*, CASE "+
			"WHEN LOWER(title) LIKE ? THEN 3 "+
			"WHEN LOWER(author) LIKE ? THEN 2 "+
			"WHEN LOWER(isbn) LIKE ? THEN 1 "+
			"ELSE 0 END as relevance_score",
			searchPattern, searchPattern, searchPattern).
		Where("LOWER(title) LIKE ? OR LOWER(author) LIKE ? OR LOWER(isbn) LIKE ?",
			searchPattern, searchPattern, searchPattern).
		Order("relevance_score DESC, title ASC").
		Find(&results).Error

	return results, err
}

// GetBooksByCategory 根據分類獲取書籍
func (s *BookServiceGORM) GetBooksByCategory(category string) ([]models.BookGORM, error) {
	var books []models.BookGORM
	err := s.db.Where("category = ?", category).Find(&books).Error
	return books, err
}

// GetBooksByPriceRange 根據價格範圍獲取書籍
func (s *BookServiceGORM) GetBooksByPriceRange(minPrice, maxPrice float64) ([]models.BookGORM, error) {
	var books []models.BookGORM
	err := s.db.Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&books).Error
	return books, err
}

// GetBooksByPublishedDate 根據出版日期獲取書籍
func (s *BookServiceGORM) GetBooksByPublishedDate(year int) ([]models.BookGORM, error) {
	var books []models.BookGORM
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	err := s.db.Where("published BETWEEN ? AND ?", startDate, endDate).Find(&books).Error
	return books, err
}

// GetBookStatistics 獲取書籍統計信息
func (s *BookServiceGORM) GetBookStatistics() (*models.BookStatistics, error) {
	var stats models.BookStatistics

	// 總書籍數
	s.db.Model(&models.BookGORM{}).Count(&stats.TotalBooks)

	// 平均價格
	s.db.Model(&models.BookGORM{}).Select("AVG(price)").Scan(&stats.AveragePrice)

	// 作者總數
	s.db.Model(&models.BookGORM{}).Distinct("author").Count(&stats.TotalAuthors)

	// 最高價格
	s.db.Model(&models.BookGORM{}).Select("MAX(price)").Scan(&stats.MostExpensive)

	// 最低價格
	s.db.Model(&models.BookGORM{}).Select("MIN(price)").Scan(&stats.LeastExpensive)

	return &stats, nil
}

// GetAuthorStatistics 獲取作者統計信息
func (s *BookServiceGORM) GetAuthorStatistics() ([]models.AuthorStatistics, error) {
	var stats []models.AuthorStatistics

	err := s.db.Model(&models.BookGORM{}).
		Select("author, COUNT(*) as book_count, AVG(price) as average_price, SUM(price) as total_sales").
		Group("author").
		Order("book_count DESC").
		Find(&stats).Error

	return stats, err
}

// GetTopRatedBooks 獲取評分最高的書籍（如果有評論功能）
func (s *BookServiceGORM) GetTopRatedBooks(limit int) ([]models.BookGORM, error) {
	var books []models.BookGORM

	err := s.db.Model(&models.BookGORM{}).
		Select("books_gorm.*, AVG(reviews_gorm.rating) as avg_rating").
		Joins("LEFT JOIN reviews_gorm ON books_gorm.id = reviews_gorm.book_id").
		Group("books_gorm.id").
		Having("AVG(reviews_gorm.rating) IS NOT NULL").
		Order("avg_rating DESC").
		Limit(limit).
		Find(&books).Error

	return books, err
}

// GetRecentBooks 獲取最近添加的書籍
func (s *BookServiceGORM) GetRecentBooks(limit int) ([]models.BookGORM, error) {
	var books []models.BookGORM
	err := s.db.Order("created_at DESC").Limit(limit).Find(&books).Error
	return books, err
}

// BulkCreateBooks 批量創建書籍
func (s *BookServiceGORM) BulkCreateBooks(books []models.BookGORM) error {
	return s.db.CreateInBatches(books, 100).Error
}

// BulkUpdateBooks 批量更新書籍
func (s *BookServiceGORM) BulkUpdateBooks(books []models.BookGORM) error {
	return s.db.Save(books).Error
}

// GetBooksWithPagination 分頁獲取書籍
func (s *BookServiceGORM) GetBooksWithPagination(page, pageSize int) ([]models.BookGORM, int64, error) {
	var books []models.BookGORM
	var total int64

	// 獲取總數
	s.db.Model(&models.BookGORM{}).Count(&total)

	// 分頁查詢
	offset := (page - 1) * pageSize
	err := s.db.Offset(offset).Limit(pageSize).Find(&books).Error

	return books, total, err
}

// GetBooksByMultipleAuthors 根據多個作者獲取書籍
func (s *BookServiceGORM) GetBooksByMultipleAuthors(authors []string) ([]models.BookGORM, error) {
	var books []models.BookGORM
	err := s.db.Where("author IN ?", authors).Find(&books).Error
	return books, err
}

// GetBooksWithReviews 獲取帶有評論的書籍
func (s *BookServiceGORM) GetBooksWithReviews() ([]models.BookGORM, error) {
	var books []models.BookGORM
	err := s.db.Preload("Reviews").Find(&books).Error
	return books, err
}

// CreateBookWithCategory 創建書籍並關聯分類
func (s *BookServiceGORM) CreateBookWithCategory(book models.BookGORM, categoryName string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 查找或創建分類
		var category models.CategoryGORM
		if err := tx.Where("name = ?", categoryName).FirstOrCreate(&category, models.CategoryGORM{Name: categoryName}).Error; err != nil {
			return err
		}

		// 設置分類 ID
		book.Category = categoryName

		// 創建書籍
		return tx.Create(&book).Error
	})
}

// GetBooksByUserFavorites 獲取用戶收藏的書籍
func (s *BookServiceGORM) GetBooksByUserFavorites(userID uint) ([]models.BookGORM, error) {
	var books []models.BookGORM
	err := s.db.Joins("JOIN user_books_gorm ON books_gorm.id = user_books_gorm.book_id").
		Where("user_books_gorm.user_id = ?", userID).
		Find(&books).Error
	return books, err
}

// AddBookToUserFavorites 添加書籍到用戶收藏
func (s *BookServiceGORM) AddBookToUserFavorites(userID, bookID uint) error {
	userBook := models.UserBookGORM{
		UserID: userID,
		BookID: bookID,
	}
	return s.db.Create(&userBook).Error
}

// RemoveBookFromUserFavorites 從用戶收藏中移除書籍
func (s *BookServiceGORM) RemoveBookFromUserFavorites(userID, bookID uint) error {
	return s.db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&models.UserBookGORM{}).Error
}

// GetDatabaseHealth 獲取數據庫健康狀態
func (s *BookServiceGORM) GetDatabaseHealth() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}

	return nil
}
