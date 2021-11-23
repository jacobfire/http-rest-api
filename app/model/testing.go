package model

import "testing"

// TestUser creates mock user
func TestUser(t *testing.T) *User {
	return &User {
		Email: "user@example.org",
		Password: "password",
	}
}
