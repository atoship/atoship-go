package atoship

import "context"

// AccountService reads the account behind the API key.
type AccountService struct{ client *Client }

// Organization is the account's organization.
type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Balance is the wallet labels are bought from.
//
// Available is what can be spent now. Pending is money committed to purchases
// that have not settled — it is not spendable, so check Available rather than
// their sum before a batch.
type Balance struct {
	Available float64 `json:"available"`
	Pending   float64 `json:"pending"`
	Currency  string  `json:"currency"`
}

// Credit is the credit line, where one is extended.
type Credit struct {
	Limit     float64 `json:"limit"`
	Used      float64 `json:"used"`
	Remaining float64 `json:"remaining"`
}

// AutoRecharge tops the wallet up when it falls below Threshold.
//
// Worth checking before a large batch: with it disabled, a batch that outruns
// the balance stops partway, and the shipments after that point come back
// unbought rather than queued.
type AutoRecharge struct {
	Enabled   bool    `json:"enabled"`
	Amount    float64 `json:"amount"`
	Threshold float64 `json:"threshold"`
}

// Account is the account the API key belongs to.
type Account struct {
	ID           string       `json:"id"`
	Object       string       `json:"object"`
	Email        string       `json:"email"`
	Name         string       `json:"name"`
	Organization Organization `json:"organization"`
	Balance      Balance      `json:"balance"`
	Credit       Credit       `json:"credit"`
	AutoRecharge AutoRecharge `json:"auto_recharge"`
}

// Get returns the account, including the spendable balance.
//
// It is also the cheapest way to check that a key works: it takes no arguments
// and buys nothing.
func (s *AccountService) Get(ctx context.Context) (*Account, error) {
	var out Account
	if err := s.client.get(ctx, "/api/v1/account", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
