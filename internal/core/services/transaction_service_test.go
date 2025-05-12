package services

import (
	"errors"
	"testing"

	"cushon/internal/core/domain"

	"github.com/shopspring/decimal"
)

// MockTransactionRepository implements output.TransactionRepository for testing
type MockTransactionRepository struct {
	transactions map[string]*domain.Transaction
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions: make(map[string]*domain.Transaction),
	}
}

func (m *MockTransactionRepository) Save(transaction *domain.Transaction) error {
	m.transactions[transaction.ID] = transaction
	return nil
}

func (m *MockTransactionRepository) FindByID(id string) (*domain.Transaction, error) {
	if transaction, exists := m.transactions[id]; exists {
		return transaction, nil
	}
	return nil, errors.New("transaction not found")
}

func (m *MockTransactionRepository) FindByUserID(userID string) ([]*domain.Transaction, error) {
	var userTransactions []*domain.Transaction
	for _, transaction := range m.transactions {
		if transaction.UserID == userID {
			userTransactions = append(userTransactions, transaction)
		}
	}
	return userTransactions, nil
}

func (m *MockTransactionRepository) Update(transaction *domain.Transaction) error {
	if _, exists := m.transactions[transaction.ID]; !exists {
		return errors.New("transaction not found")
	}
	m.transactions[transaction.ID] = transaction
	return nil
}

func (m *MockTransactionRepository) Delete(id string) error {
	if _, exists := m.transactions[id]; !exists {
		return errors.New("transaction not found")
	}
	delete(m.transactions, id)
	return nil
}

func TestTransactionService_CreateTransaction(t *testing.T) {
	repo := NewMockTransactionRepository()
	service := NewTransactionService(repo)

	tests := []struct {
		name          string
		userID        string
		amount        decimal.Decimal
		fundName      domain.FundName
		expectedError bool
	}{
		{
			name:          "valid transaction",
			userID:        "user123",
			amount:        decimal.NewFromFloat(1000.50),
			fundName:      "Cushon Equities Fund",
			expectedError: false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			amount:        decimal.NewFromFloat(1000.50),
			fundName:      "Cushon Equities Fund",
			expectedError: true,
		},
		{
			name:          "zero amount",
			userID:        "user123",
			amount:        decimal.Zero,
			fundName:      "Cushon Equities Fund",
			expectedError: true,
		},
		{
			name:          "invalid fund name",
			userID:        "user123",
			amount:        decimal.NewFromFloat(1000.50),
			fundName:      "Invalid Fund",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := service.CreateTransaction(tt.userID, tt.amount, tt.fundName)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if transaction.UserID != tt.userID {
				t.Errorf("Expected UserID %s, got %s", tt.userID, transaction.UserID)
			}

			if !transaction.Amount.Equal(tt.amount) {
				t.Errorf("Expected Amount %s, got %s", tt.amount.String(), transaction.Amount.String())
			}

			if transaction.FundName != tt.fundName {
				t.Errorf("Expected FundName %s, got %s", tt.fundName, transaction.FundName)
			}

			// Verify transaction was saved in repository
			savedTransaction, err := repo.FindByID(transaction.ID)
			if err != nil {
				t.Errorf("Failed to find saved transaction: %v", err)
			}
			if savedTransaction.UserID != tt.userID {
				t.Errorf("Expected saved transaction UserID %s, got %s", tt.userID, savedTransaction.UserID)
			}
		})
	}
}

func TestTransactionService_GetTransaction(t *testing.T) {
	repo := NewMockTransactionRepository()
	service := NewTransactionService(repo)

	// Create a test transaction
	testTransaction, _ := service.CreateTransaction(
		"user123",
		decimal.NewFromFloat(1000.50),
		"Cushon Equities Fund",
	)

	tests := []struct {
		name          string
		transactionID string
		expectedError bool
	}{
		{
			name:          "existing transaction",
			transactionID: testTransaction.ID,
			expectedError: false,
		},
		{
			name:          "non-existent transaction",
			transactionID: "non-existent",
			expectedError: true,
		},
		{
			name:          "empty ID",
			transactionID: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := service.GetTransaction(tt.transactionID)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if transaction.ID != tt.transactionID {
				t.Errorf("Expected transaction ID %s, got %s", tt.transactionID, transaction.ID)
			}
		})
	}
}

func TestTransactionService_GetUserTransactions(t *testing.T) {
	repo := NewMockTransactionRepository()
	service := NewTransactionService(repo)

	// Create test transactions for a user
	userID := "user123"
	service.CreateTransaction(userID, decimal.NewFromFloat(1000.50), "Cushon Equities Fund")
	service.CreateTransaction(userID, decimal.NewFromFloat(2000.75), "Cushon Equities Fund")

	tests := []struct {
		name          string
		userID        string
		expectedCount int
		expectedError bool
	}{
		{
			name:          "user with transactions",
			userID:        userID,
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "user without transactions",
			userID:        "no-transactions",
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "empty user ID",
			userID:        "",
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions, err := service.GetUserTransactions(tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(transactions) != tt.expectedCount {
				t.Errorf("Expected %d transactions, got %d", tt.expectedCount, len(transactions))
			}

			for _, transaction := range transactions {
				if transaction.UserID != tt.userID {
					t.Errorf("Expected UserID %s, got %s", tt.userID, transaction.UserID)
				}
			}
		})
	}
}

func TestTransactionService_UpdateTransaction(t *testing.T) {
	repo := NewMockTransactionRepository()
	service := NewTransactionService(repo)

	// Create a test transaction
	testTransaction, _ := service.CreateTransaction(
		"user123",
		decimal.NewFromFloat(1000.50),
		"Cushon Equities Fund",
	)

	tests := []struct {
		name          string
		transaction   *domain.Transaction
		expectedError bool
	}{
		{
			name: "valid update",
			transaction: &domain.Transaction{
				ID:       testTransaction.ID,
				UserID:   testTransaction.UserID,
				Amount:   decimal.NewFromFloat(2000.75),
				FundName: "Cushon Equities Fund",
			},
			expectedError: false,
		},
		{
			name: "non-existent transaction",
			transaction: &domain.Transaction{
				ID:       "non-existent",
				UserID:   "user123",
				Amount:   decimal.NewFromFloat(2000.75),
				FundName: "Cushon Equities Fund",
			},
			expectedError: true,
		},
		{
			name: "invalid fund name",
			transaction: &domain.Transaction{
				ID:       testTransaction.ID,
				UserID:   testTransaction.UserID,
				Amount:   decimal.NewFromFloat(2000.75),
				FundName: "Invalid Fund",
			},
			expectedError: true,
		},
		{
			name:          "nil transaction",
			transaction:   nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateTransaction(tt.transaction)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify update in repository
			updatedTransaction, err := repo.FindByID(tt.transaction.ID)
			if err != nil {
				t.Errorf("Failed to find updated transaction: %v", err)
			}
			if !updatedTransaction.Amount.Equal(tt.transaction.Amount) {
				t.Errorf("Expected updated amount %s, got %s", tt.transaction.Amount.String(), updatedTransaction.Amount.String())
			}
		})
	}
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	repo := NewMockTransactionRepository()
	service := NewTransactionService(repo)

	// Create a test transaction
	testTransaction, _ := service.CreateTransaction(
		"user123",
		decimal.NewFromFloat(1000.50),
		"Cushon Equities Fund",
	)

	tests := []struct {
		name          string
		transactionID string
		expectedError bool
	}{
		{
			name:          "existing transaction",
			transactionID: testTransaction.ID,
			expectedError: false,
		},
		{
			name:          "non-existent transaction",
			transactionID: "non-existent",
			expectedError: true,
		},
		{
			name:          "empty ID",
			transactionID: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteTransaction(tt.transactionID)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify deletion from repository
			_, err = repo.FindByID(tt.transactionID)
			if err == nil {
				t.Error("Expected transaction to be deleted, but found in repository")
			}
		})
	}
} 