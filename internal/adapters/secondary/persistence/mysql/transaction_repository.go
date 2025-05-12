package mysql

import (
	"database/sql"
	"time"

	"cushon/internal/core/domain"
	"cushon/internal/core/ports/output"

	"github.com/google/uuid"
)

// TransactionRepository implements the transaction repository interface
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) output.TransactionRepository {
	return &TransactionRepository{db: db}
}

// Save persists a transaction to the database
func (r *TransactionRepository) Save(transaction *domain.Transaction) error {
	if transaction.ID == "" {
		transaction.ID = uuid.New().String()
	}

	query := `
		INSERT INTO transactions (id, user_id, amount, fund_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err := r.db.Exec(query,
		transaction.ID,
		transaction.UserID,
		transaction.Amount,
		transaction.FundName,
		now,
		now,
	)

	return err
}

// FindByID retrieves a transaction by its ID
func (r *TransactionRepository) FindByID(id string) (*domain.Transaction, error) {
	query := `
		SELECT id, user_id, amount, fund_name
		FROM transactions
		WHERE id = ?
	`

	var transaction domain.Transaction
	err := r.db.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.FundName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// FindByUserID retrieves all transactions for a user
func (r *TransactionRepository) FindByUserID(userID string) ([]*domain.Transaction, error) {
	query := `
		SELECT id, user_id, amount, fund_name
		FROM transactions
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.FundName,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

// Update modifies an existing transaction
func (r *TransactionRepository) Update(transaction *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = ?, fund_name = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		transaction.Amount,
		transaction.FundName,
		time.Now(),
		transaction.ID,
	)

	return err
}

// Delete removes a transaction from the database
func (r *TransactionRepository) Delete(id string) error {
	query := `DELETE FROM transactions WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
} 