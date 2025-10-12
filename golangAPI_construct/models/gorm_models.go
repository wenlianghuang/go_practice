package models

import (
	"time"

	"gorm.io/gorm"
)

// BookGORM GORM 版本的書籍模型
// 使用 GORM 的內建功能，包括自動時間戳和軟刪除
type BookGORM struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"size:255;not null;index" json:"title" binding:"required"`
	Author     string         `gorm:"size:255;not null;index" json:"author" binding:"required"`
	Price      float64        `gorm:"type:decimal(10,2);not null" json:"price" binding:"required,min=0"`
	ISBN       string         `gorm:"size:20;uniqueIndex" json:"isbn,omitempty"`
	Category   string         `gorm:"size:100" json:"category,omitempty"`
	CategoryID *uint          `gorm:"index" json:"category_id,omitempty"`
	Published  *time.Time     `json:"published,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// 關聯
	CategoryRef *CategoryGORM `gorm:"foreignKey:CategoryID" json:"category_ref,omitempty"`
	Reviews     []ReviewGORM  `gorm:"foreignKey:BookID" json:"reviews,omitempty"`
}

// TableName 指定表名
func (BookGORM) TableName() string {
	return "books_gorm"
}

// ToBook 轉換為原始的 Book 模型（用於向後兼容）
func (b *BookGORM) ToBook() Book {
	return Book{
		ID:        string(rune(b.ID)),
		Title:     b.Title,
		Author:    b.Author,
		Price:     b.Price,
		ISBN:      b.ISBN,
		Category:  b.Category,
		Published: b.Published,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

// FromBook 從原始 Book 模型創建 GORM 模型
func (b *BookGORM) FromBook(book Book) {
	b.Title = book.Title
	b.Author = book.Author
	b.Price = book.Price
	b.ISBN = book.ISBN
	b.Category = book.Category
	b.Published = book.Published
}

// BookPatchGORM GORM 版本的書籍部分更新模型
type BookPatchGORM struct {
	Title     *string    `json:"title,omitempty"`
	Author    *string    `json:"author,omitempty"`
	Price     *float64   `json:"price,omitempty" binding:"omitempty,min=0"`
	ISBN      *string    `json:"isbn,omitempty"`
	Category  *string    `json:"category,omitempty"`
	Published *time.Time `json:"published,omitempty"`
}

// UserGORM 用戶模型（用於擴展功能）
type UserGORM struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // 不在 JSON 中顯示密碼
	FirstName string         `gorm:"size:50" json:"first_name"`
	LastName  string         `gorm:"size:50" json:"last_name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 關聯
	Books []BookGORM `gorm:"many2many:user_books;" json:"books,omitempty"`
}

// TableName 指定表名
func (UserGORM) TableName() string {
	return "users_gorm"
}

// CategoryGORM 書籍分類模型
type CategoryGORM struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 關聯
	Books []BookGORM `gorm:"foreignKey:CategoryID" json:"books,omitempty"`
}

// TableName 指定表名
func (CategoryGORM) TableName() string {
	return "categories_gorm"
}

// ReviewGORM 書籍評論模型
type ReviewGORM struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BookID    uint           `gorm:"not null;index" json:"book_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Rating    int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string         `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 關聯
	Book BookGORM `gorm:"foreignKey:BookID" json:"book,omitempty"`
	User UserGORM `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (ReviewGORM) TableName() string {
	return "reviews_gorm"
}

// UserBookGORM 用戶收藏書籍的中間表
type UserBookGORM struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	BookID    uint      `gorm:"primaryKey" json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserBookGORM) TableName() string {
	return "user_books_gorm"
}

// BookSearchResult 書籍搜索結果
type BookSearchResult struct {
	BookGORM
	RelevanceScore float64 `json:"relevance_score,omitempty"`
}

// BookStatistics 書籍統計信息
type BookStatistics struct {
	TotalBooks     int64   `json:"total_books"`
	AveragePrice   float64 `json:"average_price"`
	TotalAuthors   int64   `json:"total_authors"`
	MostExpensive  float64 `json:"most_expensive"`
	LeastExpensive float64 `json:"least_expensive"`
}

// AuthorStatistics 作者統計信息
type AuthorStatistics struct {
	Author       string  `json:"author"`
	BookCount    int64   `json:"book_count"`
	AveragePrice float64 `json:"average_price"`
	TotalSales   float64 `json:"total_sales"`
}
