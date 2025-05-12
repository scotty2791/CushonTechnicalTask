package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Transaction represents a financial transaction in the system
type Transaction struct {
	ID       string
	UserID   string
	Amount   decimal.Decimal
	FundName FundName
}

// NewTransaction creates a new transaction instance
func NewTransaction(userID string, amount decimal.Decimal, fundName FundName) *Transaction {
	return &Transaction{
		ID:       uuid.New().String(),
		UserID:   userID,
		Amount:   amount,
		FundName: fundName,
	}
} 