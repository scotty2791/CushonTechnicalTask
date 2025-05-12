package services

import (
	"errors"

	"cushon/internal/core/domain"
)

// MockDirectUserRepository implements output.DirectUserRepository for testing
type MockDirectUserRepository struct {
	users map[string]*domain.DirectUser
}

func NewMockDirectUserRepository() *MockDirectUserRepository {
	return &MockDirectUserRepository{
		users: make(map[string]*domain.DirectUser),
	}
}

func (m *MockDirectUserRepository) Save(user *domain.DirectUser) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID == "" {
		return errors.New("user ID is required")
	}
	if user.Name == "" {
		return errors.New("name is required")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockDirectUserRepository) FindByID(id string) (*domain.DirectUser, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("direct user not found")
}

func (m *MockDirectUserRepository) FindByName(name string) (*domain.DirectUser, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	for _, user := range m.users {
		if user.Name == name {
			return user, nil
		}
	}
	return nil, errors.New("direct user not found")
}

func (m *MockDirectUserRepository) Update(user *domain.DirectUser) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID == "" {
		return errors.New("user ID is required")
	}
	if user.Name == "" {
		return errors.New("name is required")
	}
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("direct user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockDirectUserRepository) Delete(id string) error {
	if id == "" {
		return errors.New("user ID is required")
	}
	if _, exists := m.users[id]; !exists {
		return errors.New("direct user not found")
	}
	delete(m.users, id)
	return nil
} 