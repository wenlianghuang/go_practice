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

var (
	// 統一錯誤
	ErrBookNotFound = errors.New("book not found")
)

// BookServiceInterface 書籍服務接口
// 定義所有書籍相關操作的標準接口
type BookServiceInterface interface {
	GetAllBooks() []models.Book
	GetBooksByAuthor(author string) []models.Book
	GetBookByID(id string) (*models.Book, error)
	CreateBook(book models.Book) (*models.Book, error)
	UpdateBook(id string, book models.Book) (*models.Book, error)
	PatchBook(id string, patch models.BookPatch) (*models.Book, error)
	DeleteBook(id string) (*models.Book, error)
	GetBooksCount() int
}

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

// GetBooksByAuthor 根據作者獲取書籍（增強版搜索）
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

// SearchBooksAdvanced 高級搜索功能（取代傳統 DB searching）
// 支持多欄位搜索、模糊匹配、排序等功能
func (s *BookServiceGORM) SearchBooksAdvanced(query string, filters map[string]interface{}) ([]models.BookSearchResult, error) {
	var results []models.BookSearchResult

	// 構建基礎查詢
	dbQuery := s.db.Model(&models.BookGORM{})

	// 如果沒有查詢詞，返回所有書籍
	if query == "" {
		err := dbQuery.Find(&results).Error
		return results, err
	}

	searchPattern := "%" + strings.ToLower(query) + "%"

	// 構建搜索條件
	searchConditions := []string{
		"LOWER(title) LIKE ?",
		"LOWER(author) LIKE ?",
		"LOWER(isbn) LIKE ?",
		"LOWER(category) LIKE ?",
	}

	searchArgs := []interface{}{
		searchPattern, searchPattern, searchPattern, searchPattern,
	}

	// 添加額外的篩選條件
	if filters != nil {
		if minPrice, ok := filters["min_price"].(float64); ok {
			dbQuery = dbQuery.Where("price >= ?", minPrice)
		}
		if maxPrice, ok := filters["max_price"].(float64); ok {
			dbQuery = dbQuery.Where("price <= ?", maxPrice)
		}
		if category, ok := filters["category"].(string); ok && category != "" {
			dbQuery = dbQuery.Where("category = ?", category)
		}
		if year, ok := filters["year"].(int); ok && year > 0 {
			dbQuery = dbQuery.Where("EXTRACT(YEAR FROM published) = ?", year)
		}
	}

	// 執行搜索查詢
	err := dbQuery.
		Select("*, CASE "+
			"WHEN LOWER(title) LIKE ? THEN 4 "+
			"WHEN LOWER(author) LIKE ? THEN 3 "+
			"WHEN LOWER(isbn) LIKE ? THEN 2 "+
			"WHEN LOWER(category) LIKE ? THEN 1 "+
			"ELSE 0 END as relevance_score",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Where(strings.Join(searchConditions, " OR "), searchArgs...).
		Order("relevance_score DESC, title ASC").
		Find(&results).Error

	return results, err
}

// GetBooksByMultipleCriteria 多條件搜索（取代傳統 DB 的多個查詢方法）
func (s *BookServiceGORM) GetBooksByMultipleCriteria(criteria map[string]interface{}) ([]models.BookGORM, error) {
	var books []models.BookGORM

	dbQuery := s.db.Model(&models.BookGORM{})

	// 動態構建查詢條件
	if title, ok := criteria["title"].(string); ok && title != "" {
		dbQuery = dbQuery.Where("title ILIKE ?", "%"+title+"%")
	}
	if author, ok := criteria["author"].(string); ok && author != "" {
		dbQuery = dbQuery.Where("author ILIKE ?", "%"+author+"%")
	}
	if category, ok := criteria["category"].(string); ok && category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}
	if minPrice, ok := criteria["min_price"].(float64); ok {
		dbQuery = dbQuery.Where("price >= ?", minPrice)
	}
	if maxPrice, ok := criteria["max_price"].(float64); ok {
		dbQuery = dbQuery.Where("price <= ?", maxPrice)
	}
	if year, ok := criteria["year"].(int); ok && year > 0 {
		dbQuery = dbQuery.Where("EXTRACT(YEAR FROM published) = ?", year)
	}
	if isbn, ok := criteria["isbn"].(string); ok && isbn != "" {
		dbQuery = dbQuery.Where("isbn = ?", isbn)
	}

	// 排序選項
	if orderBy, ok := criteria["order_by"].(string); ok {
		switch orderBy {
		case "title":
			dbQuery = dbQuery.Order("title ASC")
		case "author":
			dbQuery = dbQuery.Order("author ASC")
		case "price_asc":
			dbQuery = dbQuery.Order("price ASC")
		case "price_desc":
			dbQuery = dbQuery.Order("price DESC")
		case "created_at":
			dbQuery = dbQuery.Order("created_at DESC")
		case "published":
			dbQuery = dbQuery.Order("published DESC")
		default:
			dbQuery = dbQuery.Order("id ASC")
		}
	} else {
		dbQuery = dbQuery.Order("id ASC")
	}

	// 分頁選項
	if limit, ok := criteria["limit"].(int); ok && limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset, ok := criteria["offset"].(int); ok && offset >= 0 {
		dbQuery = dbQuery.Offset(offset)
	}

	err := dbQuery.Find(&books).Error
	return books, err
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

	// 更新所有欄位
	bookGORM.Title = book.Title
	bookGORM.Author = book.Author
	bookGORM.Price = book.Price
	bookGORM.ISBN = book.ISBN
	bookGORM.Category = book.Category
	bookGORM.Published = book.Published

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
	if patch.ISBN != nil {
		bookGORM.ISBN = *patch.ISBN
	}
	if patch.Category != nil {
		bookGORM.Category = *patch.Category
	}
	if patch.Published != nil {
		bookGORM.Published = patch.Published
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

// DeleteBookPermanently 永久刪除書籍（硬刪除）
// 這個方法會真正從數據庫中移除記錄，無法恢復
func (s *BookServiceGORM) DeleteBookPermanently(id string) (*models.Book, error) {
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

	// 硬刪除：使用 Unscoped().Delete() 繞過軟刪除
	if err := s.db.Unscoped().Delete(&bookGORM).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

// DeleteBookWithCascade 級聯硬刪除（刪除書籍及其相關記錄）
// 這個方法會刪除書籍以及所有相關的評論記錄
func (s *BookServiceGORM) DeleteBookWithCascade(id string) (*models.Book, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, ErrBookNotFound
	}

	var bookGORM models.BookGORM
	if err := s.db.Preload("Reviews").First(&bookGORM, uint(idUint)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	book := bookGORM.ToBook()

	// 開始事務
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 先刪除相關的評論
	if len(bookGORM.Reviews) > 0 {
		if err := tx.Unscoped().Where("book_id = ?", bookGORM.ID).Delete(&models.ReviewGORM{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 再刪除書籍本身
	if err := tx.Unscoped().Delete(&bookGORM).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
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
