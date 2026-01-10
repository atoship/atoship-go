package atoship

import "context"

// UsersService handles user-related operations
type UsersService struct {
	client *Client
}

// User represents a user
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Company  string `json:"company,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Active   bool   `json:"active"`
}

// GetProfile gets the current user's profile
func (s *UsersService) GetProfile(ctx context.Context) (*User, error) {
	var user User
	err := s.client.get(ctx, "/api/profile", &user)
	return &user, err
}

// UpdateProfile updates the current user's profile
func (s *UsersService) UpdateProfile(ctx context.Context, updates *User) (*User, error) {
	var user User
	err := s.client.put(ctx, "/api/profile", updates, &user)
	return &user, err
}