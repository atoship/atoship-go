package atoship

import (
	"context"
	"net/url"
)

// CarrierAccountsService manages the carrier accounts available to your
// organization.
type CarrierAccountsService struct{ client *Client }

// Carriers you can connect your own contract for.
const (
	CarrierUSPS       = "USPS"
	CarrierFedEx      = "FEDEX"
	CarrierUPS        = "UPS"
	CarrierDHL        = "DHL"
	CarrierCanadaPost = "CANADA_POST"
)

// CarrierAccountUsage counts what has been shipped on an account.
type CarrierAccountUsage struct {
	TotalLabels   int  `json:"total_labels"`
	MonthlyLabels int  `json:"monthly_labels"`
	DailyLimit    *int `json:"daily_limit"`
	MonthlyLimit  *int `json:"monthly_limit"`
}

// CarrierAccount is a carrier account you can ship on.
//
// AccountType "byoca" is your own carrier contract, billed by the carrier
// directly to you. "system" is an account atoship provides.
type CarrierAccount struct {
	ID            string              `json:"id"`
	Object        string              `json:"object"`
	Carrier       string              `json:"carrier"`
	AccountName   string              `json:"account_name"`
	AccountNumber string              `json:"account_number"`
	AccountType   string              `json:"account_type"`
	IsActive      bool                `json:"is_active"`
	IsVerified    bool                `json:"is_verified"`
	Status        string              `json:"status"`
	Logo          string              `json:"logo"`
	Capabilities  []string            `json:"capabilities"`
	Usage         CarrierAccountUsage `json:"usage"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
}

// CarrierAccountRequest connects your own carrier contract.
//
// Carrier must be one of the Carrier* constants. Credentials are the carrier's
// own — what you would use to log into their developer portal.
type CarrierAccountRequest struct {
	Carrier       string            `json:"carrier"`
	AccountName   string            `json:"account_name,omitempty"`
	AccountNumber string            `json:"account_number,omitempty"`
	Credentials   map[string]string `json:"credentials,omitempty"`
	IsActive      *bool             `json:"is_active,omitempty"`
}

// List returns the carrier accounts available to you.
//
// Filter by one of the Carrier* constants; anything else is rejected.
func (s *CarrierAccountsService) List(ctx context.Context, carrier string) ([]CarrierAccount, error) {
	q := url.Values{}
	if carrier != "" {
		q.Set("carrier", carrier)
	}
	var out struct {
		Object string           `json:"object"`
		Data   []CarrierAccount `json:"data"`
	}
	if err := s.client.get(ctx, "/api/v1/carrier-accounts", q, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Get returns one carrier account.
func (s *CarrierAccountsService) Get(ctx context.Context, id string) (*CarrierAccount, error) {
	var out CarrierAccount
	if err := s.client.get(ctx, "/api/v1/carrier-accounts/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create connects your own carrier contract.
func (s *CarrierAccountsService) Create(ctx context.Context, req *CarrierAccountRequest) (*CarrierAccount, error) {
	var out CarrierAccount
	if err := s.client.post(ctx, "/api/v1/carrier-accounts", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes a carrier account.
func (s *CarrierAccountsService) Update(ctx context.Context, id string, req *CarrierAccountRequest) (*CarrierAccount, error) {
	var out CarrierAccount
	if err := s.client.put(ctx, "/api/v1/carrier-accounts/"+url.PathEscape(id), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete disconnects a carrier account. Labels already bought on it are
// unaffected.
func (s *CarrierAccountsService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, "/api/v1/carrier-accounts/"+url.PathEscape(id), nil)
}
