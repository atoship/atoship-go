package atoship

import "context"

// CarriersService handles carrier-related operations
type CarriersService struct {
	client *Client
}

// Carrier represents a shipping carrier
type Carrier struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Code       string   `json:"code"`
	Active     bool     `json:"active"`
	Services   []string `json:"services"`
}

// List lists all available carriers
func (s *CarriersService) List(ctx context.Context) ([]Carrier, error) {
	var carriers []Carrier
	err := s.client.get(ctx, "/api/carriers", &carriers)
	return carriers, err
}