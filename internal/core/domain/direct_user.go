// note, similar entities could be created for 'Employees' (distinct to 'Retail/Direct' users)
package domain

import (
	"github.com/google/uuid"
)

// DirectUser represents a direct user in the system
type DirectUser struct {
	ID   string
	Name string
}

// NewDirectUser creates a new direct user instance
func NewDirectUser(name string) *DirectUser {
	return &DirectUser{
		ID:   uuid.New().String(),
		Name: name,
	}
} 