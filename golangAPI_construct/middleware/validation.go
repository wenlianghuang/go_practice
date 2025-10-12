package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golangAPI_construct/responses"
)

// ValidationRule 定義單一欄位的驗證規則
// 這個結構體用於配置每個欄位的驗證要求
type ValidationRule struct {
	Field    string  // 欄位名稱，對應 JSON 中的 key
	Required bool    // 是否為必填欄位
	MinLen   int     // 字串最小長度（僅用於 string 類型）
	MaxLen   int     // 字串最大長度（僅用於 string 類型）
	Min      float64 // 數值最小值（僅用於 number 類型）
	Max      float64 // 數值最大值（僅用於 number 類型）
	Pattern  string  // 正則表達式模式（用於字串格式驗證）
	Type     string  // 數據類型："string", "number", "email", "uuid"
}

// ValidationConfig 驗證配置
// 包含一組驗證規則，用於特定 API 端點的請求驗證
type ValidationConfig struct {
	Rules []ValidationRule // 驗證規則列表
}

// RequestValidator 請求驗證中間件
// 這個中間件會：
// 1. 解析請求體中的 JSON 數據
// 2. 根據提供的規則驗證每個欄位
// 3. 將驗證後的數據存儲到 context 中供後續 handler 使用
// 4. 如果驗證失敗，返回標準化的錯誤響應
func RequestValidator(config ValidationConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 只對 POST, PUT, PATCH 請求進行驗證
			// GET 和 DELETE 請求通常不需要驗證請求體
			if r.Method == http.MethodGet || r.Method == http.MethodDelete {
				next.ServeHTTP(w, r)
				return
			}

			// 解析請求體中的 JSON 數據
			// 使用 map[string]interface{} 來處理動態的 JSON 結構
			var requestData map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
				responses.Fail(w, r, responses.NewAppError(
					http.StatusBadRequest,
					"INVALID_JSON",
					"Invalid JSON format in request body",
				))
				return
			}

			// 執行驗證邏輯
			// 如果驗證失敗，會返回 AppError 並直接響應給客戶端
			if err := validateRequest(requestData, config.Rules); err != nil {
				responses.Fail(w, r, err)
				return
			}

			// 驗證成功後，將驗證後的數據存儲到 context 中
			// 後續的 handler 可以從 context 中安全地獲取這些數據
			ctx := r.Context()
			ctx = context.WithValue(ctx, "validated_data", requestData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// validateRequest 執行實際的驗證邏輯
// 這個函數會遍歷所有驗證規則，並對每個欄位進行相應的驗證
func validateRequest(data map[string]interface{}, rules []ValidationRule) *responses.AppError {
	for _, rule := range rules {
		value, exists := data[rule.Field]

		// 檢查必填欄位
		// 如果欄位被標記為必填但不存在、為 nil 或空字串，則返回錯誤
		if rule.Required && (!exists || value == nil || value == "") {
			return responses.NewAppError(
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"Field '"+rule.Field+"' is required",
			)
		}

		// 如果欄位不存在且非必填，跳過驗證
		// 這允許可選欄位不出現在請求中
		if !exists {
			continue
		}

		// 根據配置的類型進行相應的驗證
		// 每種類型都有專門的驗證函數來處理特定的驗證邏輯
		switch rule.Type {
		case "string":
			if err := validateString(value, rule); err != nil {
				return err
			}
		case "number":
			if err := validateNumber(value, rule); err != nil {
				return err
			}
		case "email":
			if err := validateEmail(value, rule); err != nil {
				return err
			}
		case "uuid":
			if err := validateUUID(value, rule); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateString 驗證字串欄位
func validateString(value interface{}, rule ValidationRule) *responses.AppError {
	str, ok := value.(string)
	if !ok {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be a string",
		)
	}

	if rule.MinLen > 0 && len(str) < rule.MinLen {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be at least "+strconv.Itoa(rule.MinLen)+" characters long",
		)
	}

	if rule.MaxLen > 0 && len(str) > rule.MaxLen {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be no more than "+strconv.Itoa(rule.MaxLen)+" characters long",
		)
	}

	if rule.Pattern != "" {
		matched, err := regexp.MatchString(rule.Pattern, str)
		if err != nil || !matched {
			return responses.NewAppError(
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"Field '"+rule.Field+"' format is invalid",
			)
		}
	}

	return nil
}

// validateNumber 驗證數字欄位
func validateNumber(value interface{}, rule ValidationRule) *responses.AppError {
	var num float64
	var ok bool

	switch v := value.(type) {
	case float64:
		num = v
		ok = true
	case int:
		num = float64(v)
		ok = true
	case string:
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			num = parsed
			ok = true
		}
	}

	if !ok {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be a number",
		)
	}

	if rule.Min != 0 && num < rule.Min {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be at least "+strconv.FormatFloat(rule.Min, 'f', -1, 64),
		)
	}

	if rule.Max != 0 && num > rule.Max {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be no more than "+strconv.FormatFloat(rule.Max, 'f', -1, 64),
		)
	}

	return nil
}

// validateEmail 驗證電子郵件格式
func validateEmail(value interface{}, rule ValidationRule) *responses.AppError {
	str, ok := value.(string)
	if !ok {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be a string",
		)
	}

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailPattern, str)
	if err != nil || !matched {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be a valid email address",
		)
	}

	return nil
}

// validateUUID 驗證 UUID 格式
func validateUUID(value interface{}, rule ValidationRule) *responses.AppError {
	str, ok := value.(string)
	if !ok {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be a string",
		)
	}

	uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	matched, err := regexp.MatchString(uuidPattern, strings.ToLower(str))
	if err != nil || !matched {
		return responses.NewAppError(
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Field '"+rule.Field+"' must be a valid UUID",
		)
	}

	return nil
}

// 預定義的驗證規則
// 這些規則可以在路由中直接使用，避免重複定義
var (
	// BookValidationRules 書籍相關的驗證規則
	// 用於 POST /api/v1/books, PUT /api/v1/books/{id}, PATCH /api/v1/books/{id} 端點
	BookValidationRules = ValidationConfig{
		Rules: []ValidationRule{
			// 書名：必填，字串類型，長度 1-200 字元
			{Field: "title", Required: true, Type: "string", MinLen: 1, MaxLen: 200},
			// 作者：必填，字串類型，長度 1-100 字元
			{Field: "author", Required: true, Type: "string", MinLen: 1, MaxLen: 100},
			// 價格：必填，數字類型，範圍 0-10000
			{Field: "price", Required: true, Type: "number", Min: 0, Max: 10000},
			// ISBN：可選，字串類型，長度 0-20 字元
			{Field: "isbn", Required: false, Type: "string", MinLen: 0, MaxLen: 20},
			// 分類：可選，字串類型，長度 0-100 字元
			{Field: "category", Required: false, Type: "string", MinLen: 0, MaxLen: 100},
			// 出版日期：可選，ISO 8601 格式字串
			{Field: "published", Required: false, Type: "string"},
		},
	}

	// LoginValidationRules 登入相關的驗證規則
	// 用於 POST /api/v1/auth/login 端點
	LoginValidationRules = ValidationConfig{
		Rules: []ValidationRule{
			// 用戶名：必填，字串類型，長度 3-50 字元
			{Field: "username", Required: true, Type: "string", MinLen: 3, MaxLen: 50},
			// 密碼：必填，字串類型，長度 6-100 字元
			{Field: "password", Required: true, Type: "string", MinLen: 6, MaxLen: 100},
		},
	}
)

// BookPatchValidationRules PATCH 請求的驗證規則
// 所有欄位都是可選的，只驗證提供的欄位
var BookPatchValidationRules = ValidationConfig{
	Rules: []ValidationRule{
		// 書名：可選，字串類型，長度 1-200 字元
		{Field: "title", Required: false, Type: "string", MinLen: 1, MaxLen: 200},
		// 作者：可選，字串類型，長度 1-100 字元
		{Field: "author", Required: false, Type: "string", MinLen: 1, MaxLen: 100},
		// 價格：可選，數字類型，範圍 0-10000
		{Field: "price", Required: false, Type: "number", Min: 0, Max: 10000},
		// ISBN：可選，字串類型，長度 0-20 字元
		{Field: "isbn", Required: false, Type: "string", MinLen: 0, MaxLen: 20},
		// 分類：可選，字串類型，長度 0-100 字元
		{Field: "category", Required: false, Type: "string", MinLen: 0, MaxLen: 100},
		// 出版日期：可選，ISO 8601 格式字串
		{Field: "published", Required: false, Type: "string"},
	},
}
