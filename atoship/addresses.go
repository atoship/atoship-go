package atoship

import (
	"context"
	"net/url"
)

// AddressesService stores and validates addresses.
type AddressesService struct{ client *Client }

// Create saves an address to your address book so later shipments can refer to
// it by ID.
func (s *AddressesService) Create(ctx context.Context, a *Address) (*Address, error) {
	var out Address
	if err := s.client.post(ctx, "/api/v1/addresses", a, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List returns saved addresses.
func (s *AddressesService) List(ctx context.Context, opts *ListOptions) ([]Address, Pagination, error) {
	q := url.Values{}
	opts.apply(q)
	var out listEnvelope[Address]
	if err := s.client.get(ctx, "/api/v1/addresses", q, &out); err != nil {
		return nil, Pagination{}, err
	}
	return out.Data, out.Pagination, nil
}

// AddressValidation is the result of checking an address against the postal
// database.
type AddressValidation struct {
	Object      string   `json:"object"`
	Valid       bool     `json:"valid"`
	Address     Address  `json:"address"`
	Original    Address  `json:"original,omitempty"`
	Messages    []string `json:"messages,omitempty"`
	Residential *bool    `json:"residential,omitempty"`
}

// Validate checks an address and returns the postal service's corrected form.
//
// This is a paid call — each request is billed whether or not the address turns
// out to be valid.
//
// Worth doing before buying a label rather than after: a corrected ZIP or a
// missing unit number is cheap to fix here and expensive once the carrier
// applies an address-correction surcharge to a label already bought.
func (s *AddressesService) Validate(ctx context.Context, a *Address) (*AddressValidation, error) {
	var out AddressValidation
	if err := s.client.post(ctx, "/api/v1/addresses/validate", a, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Verify is an alias of Validate — same handler, same charge.
//
// Both paths exist because /verify was published before the validation endpoint
// shipped and clients were written against it. Pick one: calling both validates
// twice and is billed twice.
func (s *AddressesService) Verify(ctx context.Context, a *Address) (*AddressValidation, error) {
	var out AddressValidation
	if err := s.client.post(ctx, "/api/v1/addresses/verify", a, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
