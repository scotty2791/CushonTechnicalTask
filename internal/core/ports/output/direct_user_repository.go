package output

import "cushon/internal/core/domain"

// DirectUserRepository defines the output port for direct user persistence
type DirectUserRepository interface {
	// Save persists a direct user
	Save(user *domain.DirectUser) error
	
	// FindByID retrieves a direct user by ID
	FindByID(id string) (*domain.DirectUser, error)
	
	// Update updates an existing direct user
	Update(user *domain.DirectUser) error
	
	// Delete removes a direct user by ID
	Delete(id string) error
} 