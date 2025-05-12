package input

import (
	"github.com/shopspring/decimal"
	"cushon/internal/core/domain"
)

// TransactionService defines the input port for transaction operations
type TransactionService interface {
	// CreateTransaction creates a new transaction
	CreateTransaction(userID string, amount decimal.Decimal, fundName domain.FundName) (*domain.Transaction, error)
	
	// GetTransaction retrieves a transaction by ID
	GetTransaction(id string) (*domain.Transaction, error)
	
	// GetUserTransactions retrieves all transactions for a user
	GetUserTransactions(userID string) ([]*domain.Transaction, error)
	
	// UpdateTransaction updates an existing transaction
	UpdateTransaction(transaction *domain.Transaction) error
	
	// DeleteTransaction deletes a transaction by ID
	DeleteTransaction(id string) error
} 