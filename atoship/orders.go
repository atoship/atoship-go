package atoship

import (
	"context"
	"fmt"
	"time"
)

// OrdersService handles order-related operations
type OrdersService struct {
	client *Client
}

// Order represents an order
type Order struct {
	ID                 string           `json:"id"`
	OrderNumber        string           `json:"orderNumber"`
	Source             string           `json:"source"`
	Status             string           `json:"status"`
	RecipientName      string           `json:"recipientName"`
	RecipientCompany   string           `json:"recipientCompany,omitempty"`
	RecipientStreet1   string           `json:"recipientStreet1"`
	RecipientStreet2   string           `json:"recipientStreet2,omitempty"`
	RecipientCity      string           `json:"recipientCity"`
	RecipientState     string           `json:"recipientState"`
	RecipientPostal    string           `json:"recipientPostalCode"`
	RecipientCountry   string           `json:"recipientCountry"`
	RecipientPhone     string           `json:"recipientPhone,omitempty"`
	RecipientEmail     string           `json:"recipientEmail,omitempty"`
	SenderName         string           `json:"senderName,omitempty"`
	SenderCompany      string           `json:"senderCompany,omitempty"`
	SenderStreet1      string           `json:"senderStreet1,omitempty"`
	SenderStreet2      string           `json:"senderStreet2,omitempty"`
	SenderCity         string           `json:"senderCity,omitempty"`
	SenderState        string           `json:"senderState,omitempty"`
	SenderPostal       string           `json:"senderPostalCode,omitempty"`
	SenderCountry      string           `json:"senderCountry,omitempty"`
	SenderPhone        string           `json:"senderPhone,omitempty"`
	SenderEmail        string           `json:"senderEmail,omitempty"`
	Items              []OrderItem      `json:"items"`
	TotalWeight        float64          `json:"totalWeight"`
	WeightUnit         string           `json:"weightUnit"`
	TotalValue         float64          `json:"totalValue"`
	Currency           string           `json:"currency"`
	ShippingCost       float64          `json:"shippingCost,omitempty"`
	TrackingNumber     string           `json:"trackingNumber,omitempty"`
	CarrierService     string           `json:"carrierService,omitempty"`
	ShippedAt          *time.Time       `json:"shippedAt,omitempty"`
	DeliveredAt        *time.Time       `json:"deliveredAt,omitempty"`
	Notes              string           `json:"notes,omitempty"`
	Tags               []string         `json:"tags,omitempty"`
	Metadata           map[string]any   `json:"metadata,omitempty"`
	CreatedAt          time.Time        `json:"createdAt"`
	UpdatedAt          time.Time        `json:"updatedAt"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	Weight      float64 `json:"weight"`
	WeightUnit  string  `json:"weightUnit"`
	Description string  `json:"description,omitempty"`
	ImageURL    string  `json:"imageUrl,omitempty"`
	HSCode      string  `json:"hsCode,omitempty"`
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	OrderNumber      string           `json:"orderNumber"`
	Source           string           `json:"source,omitempty"`
	RecipientName    string           `json:"recipientName"`
	RecipientCompany string           `json:"recipientCompany,omitempty"`
	RecipientStreet1 string           `json:"recipientStreet1"`
	RecipientStreet2 string           `json:"recipientStreet2,omitempty"`
	RecipientCity    string           `json:"recipientCity"`
	RecipientState   string           `json:"recipientState"`
	RecipientPostal  string           `json:"recipientPostalCode"`
	RecipientCountry string           `json:"recipientCountry"`
	RecipientPhone   string           `json:"recipientPhone,omitempty"`
	RecipientEmail   string           `json:"recipientEmail,omitempty"`
	SenderName       string           `json:"senderName,omitempty"`
	SenderCompany    string           `json:"senderCompany,omitempty"`
	SenderStreet1    string           `json:"senderStreet1,omitempty"`
	SenderStreet2    string           `json:"senderStreet2,omitempty"`
	SenderCity       string           `json:"senderCity,omitempty"`
	SenderState      string           `json:"senderState,omitempty"`
	SenderPostal     string           `json:"senderPostalCode,omitempty"`
	SenderCountry    string           `json:"senderCountry,omitempty"`
	SenderPhone      string           `json:"senderPhone,omitempty"`
	SenderEmail      string           `json:"senderEmail,omitempty"`
	Items            []OrderItem      `json:"items"`
	WeightUnit       string           `json:"weightUnit,omitempty"`
	Currency         string           `json:"currency,omitempty"`
	Notes            string           `json:"notes,omitempty"`
	Tags             []string         `json:"tags,omitempty"`
	Metadata         map[string]any   `json:"metadata,omitempty"`
}

// UpdateOrderRequest represents a request to update an order
type UpdateOrderRequest struct {
	Status           string           `json:"status,omitempty"`
	RecipientName    string           `json:"recipientName,omitempty"`
	RecipientCompany string           `json:"recipientCompany,omitempty"`
	RecipientStreet1 string           `json:"recipientStreet1,omitempty"`
	RecipientStreet2 string           `json:"recipientStreet2,omitempty"`
	RecipientCity    string           `json:"recipientCity,omitempty"`
	RecipientState   string           `json:"recipientState,omitempty"`
	RecipientPostal  string           `json:"recipientPostalCode,omitempty"`
	RecipientCountry string           `json:"recipientCountry,omitempty"`
	RecipientPhone   string           `json:"recipientPhone,omitempty"`
	RecipientEmail   string           `json:"recipientEmail,omitempty"`
	Items            []OrderItem      `json:"items,omitempty"`
	Notes            string           `json:"notes,omitempty"`
	Tags             []string         `json:"tags,omitempty"`
	Metadata         map[string]any   `json:"metadata,omitempty"`
}

// ListOrdersOptions represents options for listing orders
type ListOrdersOptions struct {
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Status     string `json:"status,omitempty"`
	Source     string `json:"source,omitempty"`
	Search     string `json:"search,omitempty"`
	StartDate  string `json:"startDate,omitempty"`
	EndDate    string `json:"endDate,omitempty"`
	SortBy     string `json:"sortBy,omitempty"`
	SortOrder  string `json:"sortOrder,omitempty"`
}

// OrderListResponse represents a paginated list of orders
type OrderListResponse struct {
	Orders  []Order `json:"orders"`
	Total   int     `json:"total"`
	Page    int     `json:"page"`
	Limit   int     `json:"limit"`
	HasMore bool    `json:"hasMore"`
}

// Create creates a new order
func (s *OrdersService) Create(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
	var order Order
	err := s.client.post(ctx, "/api/orders", req, &order)
	return &order, err
}

// Get retrieves an order by ID
func (s *OrdersService) Get(ctx context.Context, orderID string) (*Order, error) {
	var order Order
	err := s.client.get(ctx, fmt.Sprintf("/api/orders/%s", orderID), &order)
	return &order, err
}

// Update updates an existing order
func (s *OrdersService) Update(ctx context.Context, orderID string, req *UpdateOrderRequest) (*Order, error) {
	var order Order
	err := s.client.put(ctx, fmt.Sprintf("/api/orders/%s", orderID), req, &order)
	return &order, err
}

// List lists orders with optional filters
func (s *OrdersService) List(ctx context.Context, opts *ListOrdersOptions) (*OrderListResponse, error) {
	var resp OrderListResponse
	path := "/api/orders"
	
	// Add query parameters
	if opts != nil {
		// TODO: Build query string from options
	}
	
	err := s.client.get(ctx, path, &resp)
	return &resp, err
}

// Delete deletes an order
func (s *OrdersService) Delete(ctx context.Context, orderID string) error {
	return s.client.delete(ctx, fmt.Sprintf("/api/orders/%s", orderID))
}

// Ship marks an order as shipped
func (s *OrdersService) Ship(ctx context.Context, orderID string, trackingNumber string, carrier string) (*Order, error) {
	req := map[string]string{
		"trackingNumber": trackingNumber,
		"carrier":        carrier,
	}
	var order Order
	err := s.client.post(ctx, fmt.Sprintf("/api/orders/%s/ship", orderID), req, &order)
	return &order, err
}

// Cancel cancels an order
func (s *OrdersService) Cancel(ctx context.Context, orderID string, reason string) (*Order, error) {
	req := map[string]string{
		"reason": reason,
	}
	var order Order
	err := s.client.post(ctx, fmt.Sprintf("/api/orders/%s/cancel", orderID), req, &order)
	return &order, err
}

// BulkCreate creates multiple orders in batch
func (s *OrdersService) BulkCreate(ctx context.Context, orders []*CreateOrderRequest) (*BulkCreateResponse, error) {
	var resp BulkCreateResponse
	err := s.client.post(ctx, "/api/orders/batch", map[string]interface{}{
		"orders": orders,
	}, &resp)
	return &resp, err
}

// BulkCreateResponse represents the response from bulk order creation
type BulkCreateResponse struct {
	Successful []Order        `json:"successful"`
	Failed     []FailedOrder  `json:"failed"`
}

// FailedOrder represents a failed order in bulk creation
type FailedOrder struct {
	Order CreateOrderRequest `json:"order"`
	Error string             `json:"error"`
}