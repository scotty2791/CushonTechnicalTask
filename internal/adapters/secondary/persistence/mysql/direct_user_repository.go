package mysql

import (
	"database/sql"
	"errors"

	"cushon/internal/core/domain"
	"cushon/internal/core/ports/output"
)

// DirectUserRepository implements the output.DirectUserRepository interface using MySQL
type DirectUserRepository struct {
	db *sql.DB
}

// NewDirectUserRepository creates a new MySQL direct user repository
func NewDirectUserRepository(db *sql.DB) output.DirectUserRepository {
	return &DirectUserRepository{
		db: db,
	}
}

// Save persists a direct user to the database
func (r *DirectUserRepository) Save(user *domain.DirectUser) error {
	// Check if user already exists
	existingUser, err := r.FindByID(user.ID)
	if err == nil && existingUser != nil {
		return errors.New("direct user already exists")
	}

	query := `
		INSERT INTO direct_users (id, name)
		VALUES (?, ?)
	`
	_, err = r.db.Exec(query, user.ID, user.Name)
	return err
}

// FindByID retrieves a direct user by ID
func (r *DirectUserRepository) FindByID(id string) (*domain.DirectUser, error) {
	query := `
		SELECT id, name
		FROM direct_users
		WHERE id = ?
	`
	user := &domain.DirectUser{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Update updates an existing direct user
func (r *DirectUserRepository) Update(user *domain.DirectUser) error {
	query := `
		UPDATE direct_users
		SET name = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, user.Name, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("direct user not found")
	}

	return nil
}

// Delete removes a direct user by ID
func (r *DirectUserRepository) Delete(id string) error {
	query := `
		DELETE FROM direct_users
		WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	return err
} 