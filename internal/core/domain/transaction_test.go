package domain

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewTransaction(t *testing.T) {
	userID := "user123"
	amount := decimal.NewFromFloat(1000.50)
	fundName := FundName("Cushon Equities Fund")

	transaction := NewTransaction(userID, amount, fundName)

	// Test user ID is set correctly
	if transaction.UserID != userID {
		t.Errorf("Expected UserID to be %s, got %s", userID, transaction.UserID)
	}

	// Test amount is set correctly
	if !transaction.Amount.Equal(amount) {
		t.Errorf("Expected Amount to be %s, got %s", amount.String(), transaction.Amount.String())
	}

	// Test fund name is set correctly
	if transaction.FundName != fundName {
		t.Errorf("Expected FundName to be %s, got %s", fundName, transaction.FundName)
	}

	// Test ID is generated and not empty
	if transaction.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
} 