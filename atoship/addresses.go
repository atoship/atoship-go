package atoship

import (
	"context"
	"fmt"
)

// AddressesService handles address-related operations
type AddressesService struct {
	client *Client
}

// Address represents a shipping address
type Address struct {
	ID           string `json:"id,omitempty"`
	Type         string `json:"type,omitempty"` // SHIPPING or BILLING
	Name         string `json:"name"`
	Company      string `json:"company,omitempty"`
	Street1      string `json:"street1"`
	Street2      string `json:"street2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
	Country      string `json:"country"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	IsResidential bool  `json:"isResidential,omitempty"`
	Validated    bool   `json:"validated,omitempty"`
}

// ValidateAddressRequest represents a request to validate an address
type ValidateAddressRequest struct {
	Name       string `json:"name"`
	Company    string `json:"company,omitempty"`
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// ValidateAddressResponse represents the response from address validation
type ValidateAddressResponse struct {
	IsValid     bool      `json:"isValid"`
	Address     *Address  `json:"address,omitempty"`
	Suggestions []Address `json:"suggestions,omitempty"`
	Errors      []string  `json:"errors,omitempty"`
}

// Create creates a new address
func (s *AddressesService) Create(ctx context.Context, address *Address) (*Address, error) {
	var result Address
	err := s.client.post(ctx, "/api/addresses", address, &result)
	return &result, err
}

// Get retrieves an address by ID
func (s *AddressesService) Get(ctx context.Context, addressID string) (*Address, error) {
	var address Address
	err := s.client.get(ctx, fmt.Sprintf("/api/addresses/%s", addressID), &address)
	return &address, err
}

// Update updates an existing address
func (s *AddressesService) Update(ctx context.Context, addressID string, address *Address) (*Address, error) {
	var result Address
	err := s.client.put(ctx, fmt.Sprintf("/api/addresses/%s", addressID), address, &result)
	return &result, err
}

// List lists all addresses
func (s *AddressesService) List(ctx context.Context) ([]Address, error) {
	var addresses []Address
	err := s.client.get(ctx, "/api/addresses", &addresses)
	return addresses, err
}

// Delete deletes an address
func (s *AddressesService) Delete(ctx context.Context, addressID string) error {
	return s.client.delete(ctx, fmt.Sprintf("/api/addresses/%s", addressID))
}

// Validate validates an address
func (s *AddressesService) Validate(ctx context.Context, req *ValidateAddressRequest) (*ValidateAddressResponse, error) {
	var resp ValidateAddressResponse
	err := s.client.post(ctx, "/api/addresses/validate", req, &resp)
	return &resp, err
}

// Search searches for addresses with autocomplete
func (s *AddressesService) Search(ctx context.Context, query string, country string) ([]Address, error) {
	var addresses []Address
	path := fmt.Sprintf("/api/address-search?q=%s", query)
	if country != "" {
		path += fmt.Sprintf("&country=%s", country)
	}
	err := s.client.get(ctx, path, &addresses)
	return addresses, err
}