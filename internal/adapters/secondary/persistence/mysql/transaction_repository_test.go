package mysql

import (
	"database/sql"
	"testing"

	"cushon/internal/core/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func setupTransactionTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *TransactionRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	repo := NewTransactionRepository(db).(*TransactionRepository)
	return db, mock, repo
}

func TestTransactionRepository_Save(t *testing.T) {
	db, mock, repo := setupTransactionTestDB(t)
	defer db.Close()

	amount := decimal.NewFromFloat(25000.0)
	transaction := domain.NewTransaction("user123", amount, domain.FundName("Cushon Equities Fund"))
	expectedID := transaction.ID
	expectedUserID := transaction.UserID
	expectedAmount := transaction.Amount.String()
	expectedFundName := string(transaction.FundName)

	mock.ExpectExec("INSERT INTO transactions").
		WithArgs(expectedID, expectedUserID, expectedAmount, expectedFundName, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Save(transaction)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_FindByID(t *testing.T) {
	db, mock, repo := setupTransactionTestDB(t)
	defer db.Close()

	expectedID := "test-id"
	expectedUserID := "user123"
	expectedAmount, err := decimal.NewFromString("25000.0000")
	assert.NoError(t, err)
	expectedFundName := "Cushon Equities Fund"

	rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "fund_name"}).
		AddRow(expectedID, expectedUserID, expectedAmount.String(), expectedFundName)

	mock.ExpectQuery("SELECT id, user_id, amount, fund_name FROM transactions").
		WithArgs(expectedID).
		WillReturnRows(rows)

	transaction, err := repo.FindByID(expectedID)
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, expectedID, transaction.ID)
	assert.Equal(t, expectedUserID, transaction.UserID)
	assert.True(t, expectedAmount.Equal(transaction.Amount))
	assert.Equal(t, domain.FundName(expectedFundName), transaction.FundName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_FindByID_NotFound(t *testing.T) {
	db, mock, repo := setupTransactionTestDB(t)
	defer db.Close()

	expectedID := "non-existent"

	mock.ExpectQuery("SELECT id, user_id, amount, fund_name FROM transactions").
		WithArgs(expectedID).
		WillReturnError(sql.ErrNoRows)

	transaction, err := repo.FindByID(expectedID)
	assert.NoError(t, err)
	assert.Nil(t, transaction)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_FindByUserID(t *testing.T) {
	db, mock, repo := setupTransactionTestDB(t)
	defer db.Close()

	expectedUserID := "user123"

	rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "fund_name"}).
		AddRow("id1", expectedUserID, "25000.0000", "Cushon Equities Fund").
		AddRow("id2", expectedUserID, "15000.0000", "Cushon Growth Fund")

	mock.ExpectQuery("SELECT id, user_id, amount, fund_name FROM transactions").
		WithArgs(expectedUserID).
		WillReturnRows(rows)

	transactions, err := repo.FindByUserID(expectedUserID)
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, expectedUserID, transactions[0].UserID)
	assert.Equal(t, expectedUserID, transactions[1].UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_FindByUserID_Empty(t *testing.T) {
	db, mock, repo := setupTransactionTestDB(t)
	defer db.Close()

	expectedUserID := "user123"

	rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "fund_name"})

	mock.ExpectQuery("SELECT id, user_id, amount, fund_name FROM transactions").
		WithArgs(expectedUserID).
		WillReturnRows(rows)

	transactions, err := repo.FindByUserID(expectedUserID)
	assert.NoError(t, err)
	assert.Empty(t, transactions)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_Update(t *testing.T) {
	db, mock, repo := setupTransactionTestDB(t)
	defer db.Close()

	amount := decimal.NewFromFloat(25000.0)
	transaction := domain.NewTransaction("user123", amount, domain.FundName("Cushon Equities Fund"))
	transaction.FundName = domain.FundName("Cushon Growth Fund")
	expectedID := transaction.ID
	expectedAmount := transaction.Amount.String()
	expectedFundName := string(transaction.FundName)

	mock.ExpectExec("UPDATE transactions").
		WithArgs(expectedAmount, expectedFundName, sqlmock.AnyArg(), expectedID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(transaction)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_Delete(t *testing.T) {
	db, mock, repo := setupTransactionTestDB(t)
	defer db.Close()

	expectedID := "test-id"

	mock.ExpectExec("DELETE FROM transactions").
		WithArgs(expectedID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(expectedID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepository_Delete_NotFound(t *testing.T) {
	db, mock, repo := setupTransactionTestDB(t)
	defer db.Close()

	expectedID := "non-existent"

	mock.ExpectExec("DELETE FROM transactions").
		WithArgs(expectedID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(expectedID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
} 