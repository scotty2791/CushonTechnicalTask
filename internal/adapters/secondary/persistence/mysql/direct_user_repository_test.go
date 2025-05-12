package mysql

import (
	"database/sql"
	"testing"

	"cushon/internal/core/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupDirectUserTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *DirectUserRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	repo := NewDirectUserRepository(db).(*DirectUserRepository)
	return db, mock, repo
}

func TestDirectUserRepository_Save(t *testing.T) {
	db, mock, repo := setupDirectUserTestDB(t)
	defer db.Close()

	user := domain.NewDirectUser("John Doe")
	expectedID := user.ID
	expectedName := user.Name

	mock.ExpectExec("INSERT INTO direct_users").
		WithArgs(expectedID, expectedName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Save(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDirectUserRepository_FindByID(t *testing.T) {
	db, mock, repo := setupDirectUserTestDB(t)
	defer db.Close()

	expectedID := "test-id"
	expectedName := "John Doe"

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedID, expectedName)

	mock.ExpectQuery("SELECT id, name FROM direct_users").
		WithArgs(expectedID).
		WillReturnRows(rows)

	user, err := repo.FindByID(expectedID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedID, user.ID)
	assert.Equal(t, expectedName, user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDirectUserRepository_FindByID_NotFound(t *testing.T) {
	db, mock, repo := setupDirectUserTestDB(t)
	defer db.Close()

	expectedID := "non-existent"

	mock.ExpectQuery("SELECT id, name FROM direct_users").
		WithArgs(expectedID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByID(expectedID)
	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDirectUserRepository_Update(t *testing.T) {
	db, mock, repo := setupDirectUserTestDB(t)
	defer db.Close()

	user := domain.NewDirectUser("John Doe")
	user.Name = "Jane Doe"
	expectedID := user.ID
	expectedName := user.Name

	mock.ExpectExec("UPDATE direct_users").
		WithArgs(expectedName, expectedID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDirectUserRepository_Delete(t *testing.T) {
	db, mock, repo := setupDirectUserTestDB(t)
	defer db.Close()

	expectedID := "test-id"

	mock.ExpectExec("DELETE FROM direct_users").
		WithArgs(expectedID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(expectedID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDirectUserRepository_Delete_NotFound(t *testing.T) {
	db, mock, repo := setupDirectUserTestDB(t)
	defer db.Close()

	expectedID := "non-existent"

	mock.ExpectExec("DELETE FROM direct_users").
		WithArgs(expectedID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(expectedID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
} 