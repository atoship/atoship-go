package atoship

import (
	"context"
	"net/url"
	"strings"
)

// OrdersService reads orders synced from a connected store.
type OrdersService struct{ client *Client }

// Order is an order atoship synced from a connected store.
type Order struct {
	CustomerEmail   string   `json:"customer_email"`
	CustomerPhone   string   `json:"customer_phone"`
	ShippingAddress *Address `json:"shipping_address"`
	BillingAddress  *Address `json:"billing_address"`
	ShipFromAddress *Address `json:"ship_from_address"`
	Weight          float64  `json:"weight"`
	WeightUnit      string   `json:"weight_unit"`
	Length          float64  `json:"length"`
	Width           float64  `json:"width"`
	Height          float64  `json:"height"`
	DimensionUnit   string   `json:"dimension_unit"`
	PackageType     string   `json:"package_type"`
	Items           []any    `json:"items"`
	TotalAmount     float64  `json:"total_amount"`
	Currency        string   `json:"currency"`
	Tags            []string `json:"tags"`
	CreatedAt       string   `json:"created_at"`
}

// GetByPlatform looks up a synced order by the store's own order ID.
//
// This exists because a platform will not always hand its own app the customer
// address at request time — Shopify's Protected Customer Data rules can return a
// null shipping address from a live query. atoship received the address through
// the webhook when the order was created, so it can be served from here instead.
//
// Both raw IDs and Shopify GIDs are accepted; the gid://shopify/Order/ prefix is
// stripped for you.
func (s *OrdersService) GetByPlatform(ctx context.Context, platform, platformOrderID string) (*Order, error) {
	id := platformOrderID
	if i := strings.LastIndex(id, "/"); i >= 0 && strings.HasPrefix(id, "gid://") {
		id = id[i+1:]
	}
	var out struct {
		Object string `json:"object"`
		Order  Order  `json:"order"`
	}
	path := "/api/v1/orders/by-platform/" + url.PathEscape(platform) + "/" + url.PathEscape(id)
	if err := s.client.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out.Order, nil
}
