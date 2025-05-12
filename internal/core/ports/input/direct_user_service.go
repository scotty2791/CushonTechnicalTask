package input

import "cushon/internal/core/domain"

// DirectUserService defines the input port for direct user operations
type DirectUserService interface {
	// CreateDirectUser creates a new direct user
	CreateDirectUser(name string) (*domain.DirectUser, error)
	
	// GetDirectUser retrieves a direct user by ID
	GetDirectUser(id string) (*domain.DirectUser, error)
	
	// UpdateDirectUser updates an existing direct user
	UpdateDirectUser(user *domain.DirectUser) error
	
	// DeleteDirectUser deletes a direct user by ID
	DeleteDirectUser(id string) error
} 