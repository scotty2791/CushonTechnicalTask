package domain

import (
	"testing"
)

func TestNewDirectUser(t *testing.T) {
	name := "John Doe"
	user := NewDirectUser(name)

	// Test name is set correctly
	if user.Name != name {
		t.Errorf("Expected name to be %s, got %s", name, user.Name)
	}

	// Test ID is generated and not empty
	if user.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}
} 