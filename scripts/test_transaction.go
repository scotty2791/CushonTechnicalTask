package main

import (
	"fmt"
	"log"
	"os"

	"cushon/internal/adapters/secondary/persistence/mysql"
	"cushon/internal/core/domain"

	"github.com/shopspring/decimal"
)

func main() {
	// Database configuration
	dbConfig := mysql.Config{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: os.Getenv("DB_PASSWORD"),
		Database: "cushon",
	}

	// Connect to database
	db, err := mysql.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create repositories
	userRepo := mysql.NewDirectUserRepository(db)
	transactionRepo := mysql.NewTransactionRepository(db)

	// Create test user
	testUser := domain.NewDirectUser("John Doe")
	if err := userRepo.Save(testUser); err != nil {
		log.Fatalf("Failed to create test user: %v", err)
	}
	fmt.Printf("Created test user with ID: %s\n", testUser.ID)

	// Create initial transaction
	initialAmount := decimal.NewFromFloat(25000.00)
	transaction := domain.NewTransaction(testUser.ID, initialAmount, domain.FundName("Cushon Equities Fund"))
	if err := transactionRepo.Save(transaction); err != nil {
		log.Fatalf("Failed to create transaction: %v", err)
	}
	fmt.Printf("Created transaction with ID: %s and amount: %s\n", transaction.ID, transaction.Amount.String())
	
	// Verify the request
	updatedTransaction, err := transactionRepo.FindByID(transaction.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve updated transaction: %v", err)
	}
	fmt.Printf("Verified transaction amount: %s\n", updatedTransaction.Amount.String())
} 