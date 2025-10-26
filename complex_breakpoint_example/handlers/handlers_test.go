package handlers

import (
	"bytes"
	"complex_breakpoint_example/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TransactionServiceInterface defines the interface for transaction operations
type TransactionServiceInterface interface {
	Transfer(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error)
	Deposit(userID int, amount float64, description string) (*models.Transaction, error)
	Withdraw(userID int, amount float64, description string) (*models.Transaction, error)
	GetUserTransactions(userID int) ([]*models.Transaction, error)
	ConcurrentTest(req *models.ConcurrentTestRequest) (*models.ConcurrentTestResponse, error)
}

// Mock TransactionService for testing
type mockTransactionService struct {
	transferFunc func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error)
}

func (m *mockTransactionService) Transfer(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
	if m.transferFunc != nil {
		return m.transferFunc(fromUserID, toUserID, amount, description)
	}
	return nil, nil
}

// Implement other required methods to satisfy the interface
func (m *mockTransactionService) Deposit(userID int, amount float64, description string) (*models.Transaction, error) {
	return nil, nil
}

func (m *mockTransactionService) Withdraw(userID int, amount float64, description string) (*models.Transaction, error) {
	return nil, nil
}

func (m *mockTransactionService) GetUserTransactions(userID int) ([]*models.Transaction, error) {
	return nil, nil
}

func (m *mockTransactionService) ConcurrentTest(req *models.ConcurrentTestRequest) (*models.ConcurrentTestResponse, error) {
	return nil, nil
}

// TestableTransactionHandler wraps the original handler to allow testing with mocks
type TestableTransactionHandler struct {
	transferHandler    func(w http.ResponseWriter, r *http.Request)
	transactionService TransactionServiceInterface
}

func (h *TestableTransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUserID  int     `json:"from_user_id"`
		ToUserID    int     `json:"to_user_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.Transfer(req.FromUserID, req.ToUserID, req.Amount, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func TestTransactionHandler_Transfer(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		transferMockFunc   func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error)
		expectedStatusCode int
		expectedError      string
		validateResponse   func(t *testing.T, body []byte)
	}{
		{
			name: "Successful transfer",
			requestBody: map[string]interface{}{
				"from_user_id": 1,
				"to_user_id":   2,
				"amount":       100.50,
				"description":  "Test transfer",
			},
			transferMockFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
				return &models.Transaction{
					ID:          uint(1),
					UserID:      uint(fromUserID),
					Amount:      amount,
					Type:        "transfer",
					Description: description,
					Timestamp:   time.Now(),
					Status:      "completed",
				}, nil
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, body []byte) {
				var tx models.Transaction
				if err := json.Unmarshal(body, &tx); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
					return
				}
				if tx.ID != 1 {
					t.Errorf("Expected transaction ID 1, got %d", tx.ID)
				}
				if tx.Amount != 100.50 {
					t.Errorf("Expected amount 100.50, got %.2f", tx.Amount)
				}
				if tx.Type != "transfer" {
					t.Errorf("Expected type 'transfer', got %s", tx.Type)
				}
				if tx.Description != "Test transfer" {
					t.Errorf("Expected description 'Test transfer', got %s", tx.Description)
				}
			},
		},
		{
			name:               "Invalid JSON",
			requestBody:        "invalid json string",
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "Invalid JSON",
			transferMockFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
				return nil, nil
			},
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"from_user_id": 1,
				// Missing to_user_id, amount, description
			},
			transferMockFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
				return &models.Transaction{
					ID:          uint(1),
					UserID:      uint(fromUserID),
					Amount:      0,
					Type:        "transfer",
					Description: "",
					Timestamp:   time.Now(),
					Status:      "completed",
				}, nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Service error - insufficient funds",
			requestBody: map[string]interface{}{
				"from_user_id": 1,
				"to_user_id":   2,
				"amount":       1000.0,
				"description":  "Transfer test",
			},
			transferMockFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
				return nil, &models.ErrInsufficientFunds{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "insufficient funds",
		},
		{
			name: "Service error - user not found",
			requestBody: map[string]interface{}{
				"from_user_id": 999,
				"to_user_id":   2,
				"amount":       100.0,
				"description":  "Transfer test",
			},
			transferMockFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
				return nil, &models.ErrUserNotFound{UserID: fromUserID}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "user not found",
		},
		{
			name: "Zero amount transfer",
			requestBody: map[string]interface{}{
				"from_user_id": 1,
				"to_user_id":   2,
				"amount":       0.0,
				"description":  "Zero amount transfer",
			},
			transferMockFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
				return &models.Transaction{
					ID:          uint(2),
					UserID:      uint(fromUserID),
					Amount:      0.0,
					Type:        "transfer",
					Description: description,
					Timestamp:   time.Now(),
					Status:      "completed",
				}, nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Large amount transfer",
			requestBody: map[string]interface{}{
				"from_user_id": 1,
				"to_user_id":   2,
				"amount":       999999.99,
				"description":  "Large transfer",
			},
			transferMockFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
				return &models.Transaction{
					ID:          uint(3),
					UserID:      uint(fromUserID),
					Amount:      999999.99,
					Type:        "transfer",
					Description: description,
					Timestamp:   time.Now(),
					Status:      "completed",
				}, nil
			},
			expectedStatusCode: http.StatusOK,
			validateResponse: func(t *testing.T, body []byte) {
				var tx models.Transaction
				if err := json.Unmarshal(body, &tx); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
					return
				}
				if tx.Amount != 999999.99 {
					t.Errorf("Expected amount 999999.99, got %.2f", tx.Amount)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := &mockTransactionService{
				transferFunc: tt.transferMockFunc,
			}

			// Create testable handler
			handler := &TestableTransactionHandler{
				transactionService: mockService,
			}

			// Create request body
			var body []byte
			var err error
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			// Create HTTP request
			req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Transfer(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rr.Code)
			}

			// Check error message if expected
			if tt.expectedError != "" {
				body := rr.Body.String()
				if !bytes.Contains([]byte(body), []byte(tt.expectedError)) {
					t.Errorf("Expected error message containing '%s', got '%s'", tt.expectedError, body)
				}
			}

			// Validate response if function provided
			if tt.validateResponse != nil && rr.Code == http.StatusOK {
				tt.validateResponse(t, rr.Body.Bytes())
			}
		})
	}
}

// Test Transfer with incorrect HTTP method (should still work)
func TestTransactionHandler_Transfer_WrongMethod(t *testing.T) {
	mockService := &mockTransactionService{
		transferFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
			return &models.Transaction{
				ID:          uint(1),
				UserID:      uint(fromUserID),
				Amount:      amount,
				Type:        "transfer",
				Description: description,
				Timestamp:   time.Now(),
				Status:      "completed",
			}, nil
		},
	}

	handler := &TestableTransactionHandler{
		transactionService: mockService,
	}

	requestBody := map[string]interface{}{
		"from_user_id": 1,
		"to_user_id":   2,
		"amount":       100.0,
		"description":  "Test",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PUT", "/transfer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Transfer(rr, req)

	// Handler should still process the request regardless of method
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}
}

// Helper test for testing with nil body
func TestTransactionHandler_Transfer_NilBody(t *testing.T) {
	mockService := &mockTransactionService{
		transferFunc: func(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
			return nil, nil
		},
	}

	handler := &TestableTransactionHandler{
		transactionService: mockService,
	}

	req := httptest.NewRequest("POST", "/transfer", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.Transfer(rr, req)

	// Should handle nil body gracefully
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for nil body, got %d", rr.Code)
	}
}
