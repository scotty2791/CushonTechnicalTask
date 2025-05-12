package services

import (
	"errors"

	"cushon/internal/core/domain"
	"cushon/internal/core/ports/input"
	"cushon/internal/core/ports/output"

	"github.com/shopspring/decimal"
)

// TransactionService implements the input.TransactionService interface
type TransactionService struct {
	transactionRepo output.TransactionRepository
}

// NewTransactionService creates a new transaction service instance
func NewTransactionService(transactionRepo output.TransactionRepository) input.TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

// CreateTransaction implements the transaction creation use case
func (s *TransactionService) CreateTransaction(userID string, amount decimal.Decimal, fundName domain.FundName) (*domain.Transaction, error) {
	// Validate input
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if amount.IsZero() {
		return nil, errors.New("amount cannot be zero")
	}
	if !fundName.IsValid() {
		return nil, errors.New("invalid fund name")
	}

	// Create new transaction
	transaction := domain.NewTransaction(userID, amount, fundName)

	// Save transaction to repository
	if err := s.transactionRepo.Save(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransaction implements the transaction retrieval use case
func (s *TransactionService) GetTransaction(id string) (*domain.Transaction, error) {
	if id == "" {
		return nil, errors.New("transaction ID is required")
	}

	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	return transaction, nil
}

// GetUserTransactions implements the user transactions retrieval use case
func (s *TransactionService) GetUserTransactions(userID string) ([]*domain.Transaction, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	transactions, err := s.transactionRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdateTransaction implements the transaction update use case
func (s *TransactionService) UpdateTransaction(transaction *domain.Transaction) error {
	if transaction == nil {
		return errors.New("transaction cannot be nil")
	}

	if transaction.ID == "" {
		return errors.New("transaction ID is required")
	}

	if !transaction.FundName.IsValid() {
		return errors.New("invalid fund name")
	}

	// Verify transaction exists
	existingTransaction, err := s.transactionRepo.FindByID(transaction.ID)
	if err != nil {
		return err
	}

	// Update only allowed fields
	existingTransaction.Amount = transaction.Amount
	existingTransaction.FundName = transaction.FundName

	return s.transactionRepo.Update(existingTransaction)
}

// DeleteTransaction implements the transaction deletion use case
func (s *TransactionService) DeleteTransaction(id string) error {
	if id == "" {
		return errors.New("transaction ID is required")
	}

	// Verify transaction exists
	_, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.transactionRepo.Delete(id)
} 