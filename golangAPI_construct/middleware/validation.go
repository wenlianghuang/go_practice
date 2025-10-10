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

// ValidationRule 定義驗證規則
type ValidationRule struct {
	Field    string
	Required bool
	MinLen   int
	MaxLen   int
	Min      float64
	Max      float64
	Pattern  string
	Type     string // "string", "number", "email", "uuid"
}

// ValidationConfig 驗證配置
type ValidationConfig struct {
	Rules []ValidationRule
}

// RequestValidator 請求驗證中間件
func RequestValidator(config ValidationConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 只對 POST, PUT, PATCH 請求進行驗證
			if r.Method == http.MethodGet || r.Method == http.MethodDelete {
				next.ServeHTTP(w, r)
				return
			}

			// 解析請求體
			var requestData map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
				responses.Fail(w, r, responses.NewAppError(
					http.StatusBadRequest,
					"INVALID_JSON",
					"Invalid JSON format in request body",
				))
				return
			}

			// 執行驗證
			if err := validateRequest(requestData, config.Rules); err != nil {
				responses.Fail(w, r, err)
				return
			}

			// 將驗證後的數據存儲到 context 中
			ctx := r.Context()
			ctx = context.WithValue(ctx, "validated_data", requestData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// validateRequest 執行實際的驗證邏輯
func validateRequest(data map[string]interface{}, rules []ValidationRule) *responses.AppError {
	for _, rule := range rules {
		value, exists := data[rule.Field]

		// 檢查必填欄位
		if rule.Required && (!exists || value == nil || value == "") {
			return responses.NewAppError(
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"Field '"+rule.Field+"' is required",
			)
		}

		// 如果欄位不存在且非必填，跳過驗證
		if !exists {
			continue
		}

		// 根據類型進行驗證
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
var (
	// BookValidationRules 書籍相關的驗證規則
	BookValidationRules = ValidationConfig{
		Rules: []ValidationRule{
			{Field: "title", Required: true, Type: "string", MinLen: 1, MaxLen: 200},
			{Field: "author", Required: true, Type: "string", MinLen: 1, MaxLen: 100},
			{Field: "price", Required: true, Type: "number", Min: 0, Max: 10000},
		},
	}

	// LoginValidationRules 登入相關的驗證規則
	LoginValidationRules = ValidationConfig{
		Rules: []ValidationRule{
			{Field: "username", Required: true, Type: "string", MinLen: 3, MaxLen: 50},
			{Field: "password", Required: true, Type: "string", MinLen: 6, MaxLen: 100},
		},
	}
)
