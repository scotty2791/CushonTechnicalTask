package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cushon/internal/core/domain"
	"cushon/internal/core/ports/input"

	"github.com/gin-gonic/gin"
)

// MockDirectUserService implements input.DirectUserService for testing
type MockDirectUserService struct {
	users map[string]*domain.DirectUser
}

func NewMockDirectUserService() *MockDirectUserService {
	return &MockDirectUserService{
		users: make(map[string]*domain.DirectUser),
	}
}

func (m *MockDirectUserService) CreateDirectUser(name string) (*domain.DirectUser, error) {
	user := domain.NewDirectUser(name)
	m.users[user.ID] = user
	return user, nil
}

func (m *MockDirectUserService) GetDirectUser(id string) (*domain.DirectUser, error) {
	if id == "" {
		return nil, errors.New("direct user ID is required")
	}

	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("direct user not found")
}

func (m *MockDirectUserService) UpdateDirectUser(user *domain.DirectUser) error {
	if user == nil {
		return errors.New("direct user cannot be nil")
	}

	if user.ID == "" {
		return errors.New("direct user ID is required")
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

func (m *MockDirectUserService) DeleteDirectUser(id string) error {
	if id == "" {
		return errors.New("direct user ID is required")
	}

	if _, exists := m.users[id]; !exists {
		return errors.New("direct user not found")
	}

	delete(m.users, id)
	return nil
}

func setupTestRouter(service input.DirectUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewDirectUserHandler(service)
	handler.RegisterRoutes(router)
	return router
}

func TestDirectUserHandler_CreateDirectUser(t *testing.T) {
	service := NewMockDirectUserService()
	router := setupTestRouter(service)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid user",
			payload: map[string]interface{}{
				"name": "John Doe",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing name",
			payload: map[string]interface{}{
				"name": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty payload",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/direct-users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response domain.DirectUser
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.Name != tt.payload["name"] {
					t.Errorf("Expected name %s, got %s", tt.payload["name"], response.Name)
				}
			}
		})
	}
}

func TestDirectUserHandler_GetDirectUser(t *testing.T) {
	service := NewMockDirectUserService()
	router := setupTestRouter(service)

	// Create a test user
	user, _ := service.CreateDirectUser("John Doe")

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "existing user",
			userID:         user.ID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existent user",
			userID:         "non-existent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/direct-users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.DirectUser
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.ID != tt.userID {
					t.Errorf("Expected ID %s, got %s", tt.userID, response.ID)
				}
			}
		})
	}
}

func TestDirectUserHandler_UpdateDirectUser(t *testing.T) {
	service := NewMockDirectUserService()
	router := setupTestRouter(service)

	// Create a test user
	user, _ := service.CreateDirectUser("John Doe")

	tests := []struct {
		name           string
		userID         string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "valid update",
			userID: user.ID,
			payload: map[string]interface{}{
				"name": "Jane Doe",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "non-existent user",
			userID: "non-existent",
			payload: map[string]interface{}{
				"name": "Jane Doe",
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "missing name",
			userID: user.ID,
			payload: map[string]interface{}{
				"name": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/direct-users/"+tt.userID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.DirectUser
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.Name != tt.payload["name"] {
					t.Errorf("Expected name %s, got %s", tt.payload["name"], response.Name)
				}
			}
		})
	}
}

func TestDirectUserHandler_DeleteDirectUser(t *testing.T) {
	service := NewMockDirectUserService()
	router := setupTestRouter(service)

	// Create a test user
	user, _ := service.CreateDirectUser("John Doe")

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "existing user",
			userID:         user.ID,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "non-existent user",
			userID:         "non-existent",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/direct-users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
} 