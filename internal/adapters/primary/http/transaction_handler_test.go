package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cushon/internal/core/domain"
	"cushon/internal/core/ports/input"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// MockTransactionService implements input.TransactionService for testing
type MockTransactionService struct {
	transactions map[string]*domain.Transaction
}

func NewMockTransactionService() *MockTransactionService {
	return &MockTransactionService{
		transactions: make(map[string]*domain.Transaction),
	}
}

func (m *MockTransactionService) CreateTransaction(userID string, amount decimal.Decimal, fundName domain.FundName) (*domain.Transaction, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if amount.IsZero() {
		return nil, errors.New("amount cannot be zero")
	}
	if !fundName.IsValid() {
		return nil, errors.New("invalid fund name")
	}

	transaction := domain.NewTransaction(userID, amount, fundName)
	m.transactions[transaction.ID] = transaction
	return transaction, nil
}

func (m *MockTransactionService) GetTransaction(id string) (*domain.Transaction, error) {
	if id == "" {
		return nil, errors.New("transaction ID is required")
	}

	if transaction, exists := m.transactions[id]; exists {
		return transaction, nil
	}
	return nil, errors.New("transaction not found")
}

func (m *MockTransactionService) GetUserTransactions(userID string) ([]*domain.Transaction, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	var userTransactions []*domain.Transaction
	for _, transaction := range m.transactions {
		if transaction.UserID == userID {
			userTransactions = append(userTransactions, transaction)
		}
	}
	return userTransactions, nil
}

func (m *MockTransactionService) UpdateTransaction(transaction *domain.Transaction) error {
	if transaction == nil {
		return errors.New("transaction cannot be nil")
	}
	if transaction.ID == "" {
		return errors.New("transaction ID is required")
	}
	if !transaction.FundName.IsValid() {
		return errors.New("invalid fund name")
	}

	existingTransaction, exists := m.transactions[transaction.ID]
	if !exists {
		return errors.New("transaction not found")
	}

	// Update only allowed fields
	existingTransaction.Amount = transaction.Amount
	existingTransaction.FundName = transaction.FundName

	m.transactions[transaction.ID] = existingTransaction
	return nil
}

func (m *MockTransactionService) DeleteTransaction(id string) error {
	if id == "" {
		return errors.New("transaction ID is required")
	}

	if _, exists := m.transactions[id]; !exists {
		return errors.New("transaction not found")
	}

	delete(m.transactions, id)
	return nil
}

func setupTransactionTestRouter(service input.TransactionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewTransactionHandler(service)
	handler.RegisterRoutes(router)
	return router
}

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	service := NewMockTransactionService()
	router := setupTransactionTestRouter(service)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "valid transaction",
			payload: map[string]interface{}{
				"user_id":   "user123",
				"amount":    "25000.0000",
				"fund_name": "Cushon Equities Fund",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "invalid transaction",
			payload: map[string]interface{}{
				"amount":    "25000.0000",
				"fund_name": "Cushon Equities Fund",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response domain.Transaction
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.UserID != tt.payload["user_id"] {
					t.Errorf("Expected user_id %s, got %s", tt.payload["user_id"], response.UserID)
				}
			}
		})
	}
}

func TestTransactionHandler_GetTransaction(t *testing.T) {
	service := NewMockTransactionService()
	router := setupTransactionTestRouter(service)

	// Create a test transaction
	amount, _ := decimal.NewFromString("25000.0000")
	transaction, _ := service.CreateTransaction("user123", amount, domain.FundName("Cushon Equities Fund"))

	tests := []struct {
		name           string
		transactionID  string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "existing transaction",
			transactionID:  transaction.ID,
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "non-existent transaction",
			transactionID:  "non-existent",
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/transactions/"+tt.transactionID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response domain.Transaction
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.ID != tt.transactionID {
					t.Errorf("Expected ID %s, got %s", tt.transactionID, response.ID)
				}
			}
		})
	}
}

func TestTransactionHandler_GetUserTransactions(t *testing.T) {
	service := NewMockTransactionService()
	router := setupTransactionTestRouter(service)

	// Create test transactions
	amount, _ := decimal.NewFromString("25000.0000")
	service.CreateTransaction("user123", amount, domain.FundName("Cushon Equities Fund"))

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "user with transactions",
			userID:         "user123",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "user without transactions",
			userID:         "user456",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/transactions/user/"+tt.userID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response []domain.Transaction
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}
			if len(response) != tt.expectedCount {
				t.Errorf("Expected %d transactions, got %d", tt.expectedCount, len(response))
			}
		})
	}
}

func TestTransactionHandler_UpdateTransaction(t *testing.T) {
	service := NewMockTransactionService()
	router := setupTransactionTestRouter(service)

	// Create a test transaction
	amount, _ := decimal.NewFromString("25000.0000")
	transaction, _ := service.CreateTransaction("user123", amount, domain.FundName("Cushon Equities Fund"))

	// Expected amount for update
	expectedAmount, _ := decimal.NewFromString("30000.0000")

	tests := []struct {
		name           string
		transactionID  string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  bool
		expectedAmount decimal.Decimal
		expectedFund   string
	}{
		{
			name:          "valid update",
			transactionID: transaction.ID,
			payload: map[string]interface{}{
				"amount":    "30000.0000",
				"fund_name": "Cushon Equities Fund",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedAmount: expectedAmount,
			expectedFund:   "Cushon Equities Fund",
		},
		{
			name:          "non-existent transaction",
			transactionID: "non-existent",
			payload: map[string]interface{}{
				"amount":    "30000.0000",
				"fund_name": "Cushon Equities Fund",
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/transactions/"+tt.transactionID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response domain.Transaction
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.ID != tt.transactionID {
					t.Errorf("Expected ID %s, got %s", tt.transactionID, response.ID)
				}
				if !response.Amount.Equal(tt.expectedAmount) {
					t.Errorf("Expected amount %s, got %s", tt.expectedAmount.String(), response.Amount.String())
				}
				if string(response.FundName) != tt.expectedFund {
					t.Errorf("Expected fund name %s, got %s", tt.expectedFund, response.FundName)
				}
			} else {
				var errorResponse map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResponse["error"] == "" {
					t.Error("Expected error message in response, got empty string")
				}
			}
		})
	}
}

func TestTransactionHandler_DeleteTransaction(t *testing.T) {
	service := NewMockTransactionService()
	router := setupTransactionTestRouter(service)

	// Create a test transaction
	amount, _ := decimal.NewFromString("25000.0000")
	transaction, _ := service.CreateTransaction("user123", amount, domain.FundName("Cushon Equities Fund"))

	tests := []struct {
		name           string
		transactionID  string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "existing transaction",
			transactionID:  transaction.ID,
			expectedStatus: http.StatusNoContent,
			expectedError:  false,
		},
		{
			name:           "non-existent transaction",
			transactionID:  "non-existent",
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/transactions/"+tt.transactionID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
} 