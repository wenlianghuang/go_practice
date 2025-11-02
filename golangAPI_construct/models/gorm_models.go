package models

import (
	"errors"
	"fmt"
	"log"
	"strings"
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
	ISBN       string         `gorm:"size:100;uniqueIndex" json:"isbn,omitempty"`
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

// ===================
// GORM Hooks - 生命週期回調函數
// ===================

// BeforeCreate Hook - 創建前的數據驗證和預處理
// 這個 Hook 會在 GORM 執行 INSERT 操作之前自動調用
// 用於數據驗證、設置默認值、生成自動編號等
func (b *BookGORM) BeforeCreate(tx *gorm.DB) error {
	log.Printf("[GORM Hook] BeforeCreate: 開始創建書籍 '%s'", b.Title)

	// 1. 數據清理和標準化
	// 去除字串兩端的空白字符，確保數據整潔
	b.Title = strings.TrimSpace(b.Title)
	b.Author = strings.TrimSpace(b.Author)
	b.Category = strings.TrimSpace(b.Category)
	b.ISBN = strings.TrimSpace(b.ISBN)

	// 2. 數據驗證
	// 驗證標題不能為空
	if len(b.Title) == 0 {
		log.Printf("[GORM Hook] BeforeCreate: 驗證失敗 - 標題不能為空")
		return errors.New("title cannot be empty")
	}

	// 驗證作者不能為空
	if len(b.Author) == 0 {
		log.Printf("[GORM Hook] BeforeCreate: 驗證失敗 - 作者不能為空")
		return errors.New("author cannot be empty")
	}

	// 驗證價格必須為正數
	if b.Price < 0 {
		log.Printf("[GORM Hook] BeforeCreate: 驗證失敗 - 價格不能為負數 (%.2f)", b.Price)
		return errors.New("price must be positive")
	}

	// 驗證價格不能超過合理範圍
	if b.Price > 10000 {
		log.Printf("[GORM Hook] BeforeCreate: 驗證失敗 - 價格過高 (%.2f)", b.Price)
		return errors.New("price cannot exceed 10000")
	}

	// 3. 自動設置默認值
	// 如果分類為空，設置默認分類
	if b.Category == "" {
		b.Category = "General"
		log.Printf("[GORM Hook] BeforeCreate: 自動設置默認分類 'General'")
	}

	// 如果 ISBN 為空，自動生成一個
	if b.ISBN == "" {
		b.ISBN = generateISBN()
		log.Printf("[GORM Hook] BeforeCreate: 自動生成 ISBN '%s'", b.ISBN)
	}

	// 如果發布日期為空，設置為當前時間
	if b.Published == nil {
		now := time.Now()
		b.Published = &now
		log.Printf("[GORM Hook] BeforeCreate: 自動設置發布日期為當前時間")
	}

	// 4. 檢查重複數據
	// 檢查是否已存在相同標題和作者的書籍
	var existingBook BookGORM
	err := tx.Where("title = ? AND author = ?", b.Title, b.Author).First(&existingBook).Error
	if err == nil {
		log.Printf("[GORM Hook] BeforeCreate: 警告 - 已存在相同標題和作者的書籍 (ID: %d)", existingBook.ID)
		return errors.New("book with same title and author already exists")
	}

	// 檢查 ISBN 是否重複
	if b.ISBN != "" {
		err = tx.Where("isbn = ?", b.ISBN).First(&existingBook).Error
		if err == nil {
			log.Printf("[GORM Hook] BeforeCreate: 驗證失敗 - ISBN 重複 '%s'", b.ISBN)
			return errors.New("ISBN already exists")
		}
	}

	log.Printf("[GORM Hook] BeforeCreate: 驗證通過，準備創建書籍 '%s'", b.Title)
	return nil
}

// AfterCreate Hook - 創建後的後續處理
// 這個 Hook 會在 GORM 成功執行 INSERT 操作之後自動調用
// 用於更新統計數據、發送通知、記錄審計日誌等
func (b *BookGORM) AfterCreate(tx *gorm.DB) error {
	log.Printf("[GORM Hook] AfterCreate: 書籍創建成功 '%s' (ID: %d)", b.Title, b.ID)

	// 1. 更新統計數據
	// 更新書籍總數統計
	updateBookCount(tx)

	// 更新作者統計
	updateAuthorStatistics(tx, b.Author)

	// 更新分類統計
	updateCategoryStatistics(tx, b.Category)

	// 2. 自動創建分類記錄（如果不存在）
	if b.Category != "" {
		ensureCategoryExists(tx, b.Category)
	}

	// 3. 記錄審計日誌
	logAudit(tx, "BOOK_CREATED", b.ID, fmt.Sprintf("Book '%s' by %s created", b.Title, b.Author))

	// 4. 發送業務通知
	sendNotification("BOOK_CREATED", map[string]interface{}{
		"book_id": b.ID,
		"title":   b.Title,
		"author":  b.Author,
		"price":   b.Price,
		"isbn":    b.ISBN,
	})

	// 5. 更新緩存（如果使用緩存）
	invalidateCache("books")

	log.Printf("[GORM Hook] AfterCreate: 後續處理完成")
	return nil
}

// BeforeUpdate Hook - 更新前的數據驗證和預處理
// 這個 Hook 會在 GORM 執行 UPDATE 操作之前自動調用
// 用於檢查數據變化、驗證更新數據、記錄變更日誌等
func (b *BookGORM) BeforeUpdate(tx *gorm.DB) error {
	log.Printf("[GORM Hook] BeforeUpdate: 開始更新書籍 '%s' (ID: %d)", b.Title, b.ID)

	// 1. 獲取原始數據用於比較
	var originalBook BookGORM
	if err := tx.First(&originalBook, b.ID).Error; err != nil {
		log.Printf("[GORM Hook] BeforeUpdate: 無法獲取原始數據 - %v", err)
		return err
	}

	// 2. 數據清理和標準化
	b.Title = strings.TrimSpace(b.Title)
	b.Author = strings.TrimSpace(b.Author)
	b.Category = strings.TrimSpace(b.Category)
	b.ISBN = strings.TrimSpace(b.ISBN)

	// 3. 數據驗證
	if len(b.Title) == 0 {
		log.Printf("[GORM Hook] BeforeUpdate: 驗證失敗 - 標題不能為空")
		return errors.New("title cannot be empty")
	}

	if len(b.Author) == 0 {
		log.Printf("[GORM Hook] BeforeUpdate: 驗證失敗 - 作者不能為空")
		return errors.New("author cannot be empty")
	}

	if b.Price < 0 {
		log.Printf("[GORM Hook] BeforeUpdate: 驗證失敗 - 價格不能為負數")
		return errors.New("price must be positive")
	}

	if b.Price > 10000 {
		log.Printf("[GORM Hook] BeforeUpdate: 驗證失敗 - 價格過高")
		return errors.New("price cannot exceed 10000")
	}

	// 4. 檢查重要欄位變化
	changes := make(map[string]interface{})

	// 檢查價格變化
	if originalBook.Price != b.Price {
		changes["price"] = map[string]interface{}{
			"old": originalBook.Price,
			"new": b.Price,
		}
		log.Printf("[GORM Hook] BeforeUpdate: 價格變化 %.2f -> %.2f", originalBook.Price, b.Price)

		// 如果價格大幅下降，記錄促銷信息
		if b.Price < originalBook.Price*0.8 {
			log.Printf("[GORM Hook] BeforeUpdate: 檢測到大幅降價，觸發促銷通知")
			sendPromotionNotification(b.Title, b.Price, originalBook.Price)
		}
	}

	// 檢查分類變化
	if originalBook.Category != b.Category {
		changes["category"] = map[string]interface{}{
			"old": originalBook.Category,
			"new": b.Category,
		}
		log.Printf("[GORM Hook] BeforeUpdate: 分類變化 '%s' -> '%s'", originalBook.Category, b.Category)
	}

	// 檢查 ISBN 變化
	if originalBook.ISBN != b.ISBN {
		changes["isbn"] = map[string]interface{}{
			"old": originalBook.ISBN,
			"new": b.ISBN,
		}
		log.Printf("[GORM Hook] BeforeUpdate: ISBN 變化 '%s' -> '%s'", originalBook.ISBN, b.ISBN)

		// 檢查新 ISBN 是否重複
		if b.ISBN != "" {
			var existingBook BookGORM
			err := tx.Where("isbn = ? AND id != ?", b.ISBN, b.ID).First(&existingBook).Error
			if err == nil {
				log.Printf("[GORM Hook] BeforeUpdate: 驗證失敗 - ISBN 重複 '%s'", b.ISBN)
				return errors.New("ISBN already exists")
			}
		}
	}

	// 5. 記錄變更日誌
	if len(changes) > 0 {
		logChangeHistory(tx, b.ID, changes)
	}

	log.Printf("[GORM Hook] BeforeUpdate: 驗證通過，準備更新書籍")
	return nil
}

// AfterUpdate Hook - 更新後的後續處理
// 這個 Hook 會在 GORM 成功執行 UPDATE 操作之後自動調用
// 用於更新相關統計、發送通知、記錄審計日誌等
func (b *BookGORM) AfterUpdate(tx *gorm.DB) error {
	log.Printf("[GORM Hook] AfterUpdate: 書籍更新成功 '%s' (ID: %d)", b.Title, b.ID)

	// 1. 更新統計數據
	updateBookCount(tx)
	updateAuthorStatistics(tx, b.Author)
	updateCategoryStatistics(tx, b.Category)

	// 2. 記錄審計日誌
	logAudit(tx, "BOOK_UPDATED", b.ID, fmt.Sprintf("Book '%s' updated", b.Title))

	// 3. 發送更新通知
	sendNotification("BOOK_UPDATED", map[string]interface{}{
		"book_id": b.ID,
		"title":   b.Title,
		"author":  b.Author,
		"price":   b.Price,
	})

	// 4. 更新緩存
	invalidateCache("books")

	log.Printf("[GORM Hook] AfterUpdate: 後續處理完成")
	return nil
}

// BeforeDelete Hook - 刪除前的檢查和預處理
// 這個 Hook 會在 GORM 執行 DELETE 操作之前自動調用
// 用於檢查是否可以刪除、記錄刪除原因等
func (b *BookGORM) BeforeDelete(tx *gorm.DB) error {
	log.Printf("[GORM Hook] BeforeDelete: 開始刪除書籍 '%s' (ID: %d)", b.Title, b.ID)

	// 1. 檢查是否有相關評論
	var reviewCount int64
	if err := tx.Model(&ReviewGORM{}).Where("book_id = ?", b.ID).Count(&reviewCount).Error; err != nil {
		log.Printf("[GORM Hook] BeforeDelete: 檢查評論時出錯 - %v", err)
		return err
	}

	if reviewCount > 0 {
		log.Printf("[GORM Hook] BeforeDelete: 無法刪除 - 書籍有 %d 條評論", reviewCount)
		return errors.New("cannot delete book with reviews")
	}

	// 2. 檢查是否有用戶收藏
	var favoriteCount int64
	if err := tx.Model(&UserBookGORM{}).Where("book_id = ?", b.ID).Count(&favoriteCount).Error; err != nil {
		log.Printf("[GORM Hook] BeforeDelete: 檢查收藏時出錯 - %v", err)
		return err
	}

	if favoriteCount > 0 {
		log.Printf("[GORM Hook] BeforeDelete: 無法刪除 - 書籍被 %d 個用戶收藏", favoriteCount)
		return errors.New("cannot delete book with user favorites")
	}

	// 3. 檢查是否為重要書籍（價格超過 100 的書籍需要額外確認）
	if b.Price > 100 {
		log.Printf("[GORM Hook] BeforeDelete: 警告 - 正在刪除高價值書籍 (價格: %.2f)", b.Price)
		// 在實際應用中，這裡可以發送管理員確認通知
	}

	// 4. 記錄刪除原因和相關信息
	log.Printf("[GORM Hook] BeforeDelete: 刪除原因檢查完成，書籍信息:")
	log.Printf("  - 標題: %s", b.Title)
	log.Printf("  - 作者: %s", b.Author)
	log.Printf("  - 價格: %.2f", b.Price)
	log.Printf("  - ISBN: %s", b.ISBN)
	log.Printf("  - 分類: %s", b.Category)

	log.Printf("[GORM Hook] BeforeDelete: 驗證通過，準備刪除書籍")
	return nil
}

// AfterDelete Hook - 刪除後的後續處理
// 這個 Hook 會在 GORM 成功執行 DELETE 操作之後自動調用
// 用於更新統計數據、發送通知、記錄審計日誌等
func (b *BookGORM) AfterDelete(tx *gorm.DB) error {
	log.Printf("[GORM Hook] AfterDelete: 書籍刪除成功 '%s' (ID: %d)", b.Title, b.ID)

	// 1. 更新統計數據
	updateBookCount(tx)
	updateAuthorStatistics(tx, b.Author)
	updateCategoryStatistics(tx, b.Category)

	// 2. 清理相關數據
	// 清理可能存在的孤立數據
	cleanupRelatedData(tx, b.ID)

	// 3. 記錄審計日誌
	logAudit(tx, "BOOK_DELETED", b.ID, fmt.Sprintf("Book '%s' deleted", b.Title))

	// 4. 發送刪除通知
	sendNotification("BOOK_DELETED", map[string]interface{}{
		"book_id": b.ID,
		"title":   b.Title,
		"author":  b.Author,
		"price":   b.Price,
	})

	// 5. 更新緩存
	invalidateCache("books")

	log.Printf("[GORM Hook] AfterDelete: 後續處理完成")
	return nil
}

// ===================
// GORM Hooks 輔助函數
// ===================

// generateISBN 生成唯一的 ISBN 號碼
// 這是一個簡單的實現，在實際應用中可能需要更複雜的邏輯
func generateISBN() string {
	// 使用時間戳和隨機數生成 ISBN
	// 格式:  (標準 ISBN 格式)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%04d-%04d-%04d-%04d",
		timestamp%10000,
		(timestamp/10000)%10000,
		(timestamp/1000000)%10000,
		(timestamp/1000000)%1000)
}

// updateBookCount 更新書籍總數統計
// 這個函數會在書籍創建、更新、刪除後調用，維護準確的統計數據
func updateBookCount(tx *gorm.DB) {
	var count int64
	if err := tx.Model(&BookGORM{}).Count(&count).Error; err != nil {
		log.Printf("[GORM Hook Helper] updateBookCount: 更新書籍總數失敗 - %v", err)
		return
	}
	log.Printf("[GORM Hook Helper] updateBookCount: 書籍總數更新為 %d", count)

	// 在實際應用中，這裡可以更新緩存或發送統計數據到監控系統
	// 例如：更新 Redis 緩存、發送到 Prometheus 等
}

// updateAuthorStatistics 更新作者統計數據
// 統計每個作者的書籍數量、平均價格等信息
func updateAuthorStatistics(tx *gorm.DB, author string) {
	var count int64
	var avgPrice float64

	// 統計該作者的書籍數量
	if err := tx.Model(&BookGORM{}).Where("author = ?", author).Count(&count).Error; err != nil {
		log.Printf("[GORM Hook Helper] updateAuthorStatistics: 統計作者書籍數量失敗 - %v", err)
		return
	}

	// 計算該作者的平均價格
	if err := tx.Model(&BookGORM{}).Where("author = ?", author).Select("AVG(price)").Scan(&avgPrice).Error; err != nil {
		log.Printf("[GORM Hook Helper] updateAuthorStatistics: 計算平均價格失敗 - %v", err)
		return
	}

	log.Printf("[GORM Hook Helper] updateAuthorStatistics: 作者 '%s' 統計更新 - 書籍數量: %d, 平均價格: %.2f",
		author, count, avgPrice)

	// 在實際應用中，這裡可以更新作者統計表或緩存
}

// updateCategoryStatistics 更新分類統計數據
// 統計每個分類的書籍數量、平均價格等信息
func updateCategoryStatistics(tx *gorm.DB, category string) {
	var count int64
	var avgPrice float64

	// 統計該分類的書籍數量
	if err := tx.Model(&BookGORM{}).Where("category = ?", category).Count(&count).Error; err != nil {
		log.Printf("[GORM Hook Helper] updateCategoryStatistics: 統計分類書籍數量失敗 - %v", err)
		return
	}

	// 計算該分類的平均價格
	if err := tx.Model(&BookGORM{}).Where("category = ?", category).Select("AVG(price)").Scan(&avgPrice).Error; err != nil {
		log.Printf("[GORM Hook Helper] updateCategoryStatistics: 計算平均價格失敗 - %v", err)
		return
	}

	log.Printf("[GORM Hook Helper] updateCategoryStatistics: 分類 '%s' 統計更新 - 書籍數量: %d, 平均價格: %.2f",
		category, count, avgPrice)

	// 在實際應用中，這裡可以更新分類統計表或緩存
}

// ensureCategoryExists 確保分類記錄存在
// 如果指定的分類不存在，自動創建一個新的分類記錄
func ensureCategoryExists(tx *gorm.DB, categoryName string) {
	var category CategoryGORM

	// 檢查分類是否已存在
	if err := tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 分類不存在，創建新的分類記錄
			category = CategoryGORM{
				Name:        categoryName,
				Description: fmt.Sprintf("Auto-created category: %s", categoryName),
			}

			if err := tx.Create(&category).Error; err != nil {
				log.Printf("[GORM Hook Helper] ensureCategoryExists: 創建分類失敗 - %v", err)
				return
			}

			log.Printf("[GORM Hook Helper] ensureCategoryExists: 自動創建分類 '%s' (ID: %d)",
				categoryName, category.ID)
		} else {
			log.Printf("[GORM Hook Helper] ensureCategoryExists: 檢查分類時出錯 - %v", err)
		}
	}
}

// logAudit 記錄審計日誌
// 記錄所有重要的數據庫操作，用於審計和追蹤
func logAudit(tx *gorm.DB, action string, bookID uint, description string) {
	log.Printf("[GORM Hook Helper] logAudit: %s - Book ID: %d - %s", action, bookID, description)

	// 在實際應用中，這裡可以將審計日誌寫入專門的審計表
	// 例如：
	// auditLog := AuditLog{
	//     Action:      action,
	//     TableName:   "books_gorm",
	//     RecordID:    bookID,
	//     Description: description,
	//     Timestamp:   time.Now(),
	// }
	// tx.Create(&auditLog)
}

// sendNotification 發送業務通知
// 在重要操作發生時發送通知，如創建、更新、刪除書籍
func sendNotification(eventType string, data map[string]interface{}) {
	log.Printf("[GORM Hook Helper] sendNotification: 發送通知 - 事件類型: %s, 數據: %+v", eventType, data)

	// 在實際應用中，這裡可以：
	// 1. 發送郵件通知
	// 2. 發送 Slack/Teams 消息
	// 3. 推送到消息隊列
	// 4. 觸發 Webhook
	// 5. 更新實時通知系統

	// 示例：根據事件類型處理不同的通知邏輯
	switch eventType {
	case "BOOK_CREATED":
		log.Printf("📚 新書上架通知: %s", data["title"])
	case "BOOK_UPDATED":
		log.Printf("📝 書籍更新通知: %s", data["title"])
	case "BOOK_DELETED":
		log.Printf("🗑️ 書籍下架通知: %s", data["title"])
	}
}

// sendPromotionNotification 發送促銷通知
// 當檢測到價格大幅下降時，發送促銷通知
func sendPromotionNotification(title string, newPrice, oldPrice float64) {
	discountPercent := (oldPrice - newPrice) / oldPrice * 100

	log.Printf("[GORM Hook Helper] sendPromotionNotification: 促銷通知 - 書籍: %s, 原價: %.2f, 現價: %.2f, 折扣: %.1f%%",
		title, oldPrice, newPrice, discountPercent)

	// 在實際應用中，這裡可以：
	// 1. 發送促銷郵件給用戶
	// 2. 更新促銷列表
	// 3. 推送到促銷系統
	// 4. 發送推播通知
}

// logChangeHistory 記錄變更歷史
// 記錄書籍重要欄位的變更歷史，用於追蹤和審計
func logChangeHistory(tx *gorm.DB, bookID uint, changes map[string]interface{}) {
	log.Printf("[GORM Hook Helper] logChangeHistory: 記錄變更歷史 - Book ID: %d, 變更: %+v", bookID, changes)

	// 在實際應用中，這裡可以將變更歷史寫入專門的變更歷史表
	// 例如：
	// for field, change := range changes {
	//     changeLog := ChangeLog{
	//         TableName:  "books_gorm",
	//         RecordID:   bookID,
	//         FieldName:  field,
	//         OldValue:   change["old"],
	//         NewValue:   change["new"],
	//         ChangedAt:  time.Now(),
	//     }
	//     tx.Create(&changeLog)
	// }
}

// cleanupRelatedData 清理相關數據
// 在刪除書籍後，清理可能存在的孤立數據
func cleanupRelatedData(tx *gorm.DB, bookID uint) {
	log.Printf("[GORM Hook Helper] cleanupRelatedData: 清理相關數據 - Book ID: %d", bookID)

	// 在實際應用中，這裡可以清理：
	// 1. 孤立的評論記錄
	// 2. 孤立的收藏記錄
	// 3. 孤立的搜索索引
	// 4. 孤立的緩存數據

	// 示例：清理孤立的評論記錄（雖然 BeforeDelete 已經檢查過）
	// tx.Where("book_id = ?", bookID).Delete(&ReviewGORM{})

	log.Printf("[GORM Hook Helper] cleanupRelatedData: 相關數據清理完成")
}

// invalidateCache 使緩存失效
// 在數據變更後，清理相關的緩存數據
func invalidateCache(cacheKey string) {
	log.Printf("[GORM Hook Helper] invalidateCache: 使緩存失效 - Key: %s", cacheKey)

	// 在實際應用中，這裡可以：
	// 1. 清理 Redis 緩存
	// 2. 清理內存緩存
	// 3. 發送緩存失效事件
	// 4. 更新 CDN 緩存

	// 示例：清理相關的緩存鍵
	cacheKeys := []string{
		"books",
		"books_list",
		"books_stats",
		fmt.Sprintf("book_%s", cacheKey),
	}

	for _, key := range cacheKeys {
		log.Printf("[GORM Hook Helper] invalidateCache: 清理緩存鍵 '%s'", key)
		// 實際的緩存清理邏輯
	}
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
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"` // 不在 JSON 中顯示密碼
	Roles        string         `gorm:"size:120;default:'user'" json:"roles"`
	FirstName    string         `gorm:"size:50" json:"first_name"`
	LastName     string         `gorm:"size:50" json:"last_name"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time     `json:"last_login,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

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
