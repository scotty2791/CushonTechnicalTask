package services

import (
	"testing"

	"cushon/internal/core/domain"
)

func TestDirectUserService_CreateDirectUser(t *testing.T) {
	repo := NewMockDirectUserRepository()
	service := NewDirectUserService(repo)

	tests := []struct {
		name          string
		inputName     string
		expectedError bool
	}{
		{
			name:          "valid name",
			inputName:     "John Doe",
			expectedError: false,
		},
		{
			name:          "empty name",
			inputName:     "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.CreateDirectUser(tt.inputName)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user.Name != tt.inputName {
				t.Errorf("Expected name %s, got %s", tt.inputName, user.Name)
			}

			if user.ID == "" {
				t.Error("Expected ID to be generated, got empty string")
			}

			// Verify user was saved in repository
			savedUser, err := repo.FindByID(user.ID)
			if err != nil {
				t.Errorf("Failed to find saved user: %v", err)
			}
			if savedUser.Name != tt.inputName {
				t.Errorf("Expected saved user name %s, got %s", tt.inputName, savedUser.Name)
			}
		})
	}
}

func TestDirectUserService_GetDirectUser(t *testing.T) {
	repo := NewMockDirectUserRepository()
	service := NewDirectUserService(repo)

	// Create a test user
	testUser, _ := service.CreateDirectUser("John Doe")

	tests := []struct {
		name          string
		userID        string
		expectedError bool
	}{
		{
			name:          "existing user",
			userID:        testUser.ID,
			expectedError: false,
		},
		{
			name:          "non-existent user",
			userID:        "non-existent",
			expectedError: true,
		},
		{
			name:          "empty ID",
			userID:        "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.GetDirectUser(tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user.ID != tt.userID {
				t.Errorf("Expected user ID %s, got %s", tt.userID, user.ID)
			}
		})
	}
}

func TestDirectUserService_UpdateDirectUser(t *testing.T) {
	repo := NewMockDirectUserRepository()
	service := NewDirectUserService(repo)

	// Create a test user
	testUser, _ := service.CreateDirectUser("John Doe")

	tests := []struct {
		name          string
		user          *domain.DirectUser
		expectedError bool
	}{
		{
			name: "valid update",
			user: &domain.DirectUser{
				ID:   testUser.ID,
				Name: "Jane Doe",
			},
			expectedError: false,
		},
		{
			name: "non-existent user",
			user: &domain.DirectUser{
				ID:   "non-existent",
				Name: "Jane Doe",
			},
			expectedError: true,
		},
		{
			name: "empty name",
			user: &domain.DirectUser{
				ID:   testUser.ID,
				Name: "",
			},
			expectedError: true,
		},
		{
			name:          "nil user",
			user:          nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateDirectUser(tt.user)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify user was updated in repository
			updatedUser, err := repo.FindByID(tt.user.ID)
			if err != nil {
				t.Errorf("Failed to find updated user: %v", err)
			}
			if updatedUser.Name != tt.user.Name {
				t.Errorf("Expected updated user name %s, got %s", tt.user.Name, updatedUser.Name)
			}
		})
	}
}

func TestDirectUserService_DeleteDirectUser(t *testing.T) {
	repo := NewMockDirectUserRepository()
	service := NewDirectUserService(repo)

	// Create a test user
	testUser, _ := service.CreateDirectUser("John Doe")

	tests := []struct {
		name          string
		userID        string
		expectedError bool
	}{
		{
			name:          "existing user",
			userID:        testUser.ID,
			expectedError: false,
		},
		{
			name:          "non-existent user",
			userID:        "non-existent",
			expectedError: true,
		},
		{
			name:          "empty ID",
			userID:        "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteDirectUser(tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify user was deleted from repository
			_, err = repo.FindByID(tt.userID)
			if err == nil {
				t.Error("Expected user to be deleted, but found in repository")
			}
		})
	}
} 