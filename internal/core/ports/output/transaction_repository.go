package output

import "cushon/internal/core/domain"

// TransactionRepository defines the output port for transaction persistence
type TransactionRepository interface {
	// Save persists a transaction
	Save(transaction *domain.Transaction) error
	
	// FindByID retrieves a transaction by ID
	FindByID(id string) (*domain.Transaction, error)
	
	// FindByUserID retrieves all transactions for a user
	FindByUserID(userID string) ([]*domain.Transaction, error)
	
	// Update updates an existing transaction
	Update(transaction *domain.Transaction) error
	
	// Delete removes a transaction by ID
	Delete(id string) error
} 