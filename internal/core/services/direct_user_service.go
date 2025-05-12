package services

import (
	"errors"

	"cushon/internal/core/domain"
	"cushon/internal/core/ports/input"
	"cushon/internal/core/ports/output"
)

// DirectUserService implements the input.DirectUserService interface
type DirectUserService struct {
	directUserRepo output.DirectUserRepository
}

// NewDirectUserService creates a new direct user service instance
func NewDirectUserService(directUserRepo output.DirectUserRepository) input.DirectUserService {
	return &DirectUserService{
		directUserRepo: directUserRepo,
	}
}

// CreateDirectUser implements the direct user creation use case
func (s *DirectUserService) CreateDirectUser(name string) (*domain.DirectUser, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Create new direct user
	directUser := domain.NewDirectUser(name)

	// Save direct user to repository
	if err := s.directUserRepo.Save(directUser); err != nil {
		return nil, err
	}

	return directUser, nil
}

// GetDirectUser implements the direct user retrieval use case
func (s *DirectUserService) GetDirectUser(id string) (*domain.DirectUser, error) {
	if id == "" {
		return nil, errors.New("direct user ID is required")
	}

	directUser, err := s.directUserRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return directUser, nil
}

// UpdateDirectUser implements the direct user update use case
func (s *DirectUserService) UpdateDirectUser(user *domain.DirectUser) error {
	if user == nil {
		return errors.New("direct user cannot be nil")
	}

	if user.ID == "" {
		return errors.New("direct user ID is required")
	}

	if user.Name == "" {
		return errors.New("name is required")
	}

	return s.directUserRepo.Update(user)
}

// DeleteDirectUser implements the direct user deletion use case
func (s *DirectUserService) DeleteDirectUser(id string) error {
	if id == "" {
		return errors.New("direct user ID is required")
	}

	// Verify direct user exists
	_, err := s.directUserRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.directUserRepo.Delete(id)
} 