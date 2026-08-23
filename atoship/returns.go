package atoship

import (
	"context"
	"net/url"
)

// ReturnsService creates and lists return labels.
type ReturnsService struct{ client *Client }

// Return is a return shipping label.
type Return struct {
	ID             string   `json:"id"`
	Object         string   `json:"object"`
	TrackingNumber string   `json:"tracking_number"`
	Carrier        string   `json:"carrier"`
	Service        string   `json:"service"`
	ServiceCode    string   `json:"service_code"`
	Status         string   `json:"status"`
	IsReturn       bool     `json:"is_return"`
	FromAddress    Address  `json:"from_address"`
	ToAddress      Address  `json:"to_address"`
	ReturnAddress  *Address `json:"return_address"`
	Weight         float64  `json:"weight"`
	WeightUnit     string   `json:"weight_unit"`
	// Cost and RetailCost are both what the customer was charged, and are null
	// on a label created before the price was recorded — null means unknown,
	// not free.
	Cost              *float64 `json:"cost"`
	RetailCost        *float64 `json:"retail_cost"`
	Currency          string   `json:"currency"`
	LabelURL          string   `json:"label_url"`
	LabelFormat       string   `json:"label_format"`
	Reference         string   `json:"reference"`
	Reason            string   `json:"reason"`
	RMANumber         string   `json:"rma_number"`
	OriginalLabelID   string   `json:"original_label_id"`
	Notes             string   `json:"notes"`
	EstimatedDelivery string   `json:"estimated_delivery"`
	ShippedAt         string   `json:"shipped_at"`
	DeliveredAt       string   `json:"delivered_at"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

// ReturnRequest creates a return label.
//
// OriginalLabelID links the return to the outbound shipment, which is what makes
// it show up against the original order. From and To are as the parcel travels
// on the way back: From is the customer, To is your warehouse.
type ReturnRequest struct {
	OriginalLabelID string     `json:"original_label_id,omitempty"`
	FromAddress     AddressRef `json:"from_address"`
	ToAddress       AddressRef `json:"to_address"`
	Parcel          ParcelRef  `json:"parcel"`
	Carrier         string     `json:"carrier,omitempty"`
	Service         string     `json:"service,omitempty"`
	Reference       string     `json:"reference,omitempty"`
	Reason          string     `json:"reason,omitempty"`
	RMANumber       string     `json:"rma_number,omitempty"`
	Notes           string     `json:"notes,omitempty"`
}

// ReturnSummary aggregates an account's returns.
type ReturnSummary struct {
	TotalReturns int     `json:"total_returns"`
	TotalCost    float64 `json:"total_cost"`
	Currency     string  `json:"currency"`
}

// Create makes a return label.
func (s *ReturnsService) Create(ctx context.Context, req *ReturnRequest) (*Return, error) {
	var out Return
	if err := s.client.post(ctx, "/api/v1/returns", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ReturnListOptions filters a return list.
type ReturnListOptions struct {
	ListOptions
	Status  string
	Carrier string
}

// List returns your return labels, with an account-level summary.
func (s *ReturnsService) List(ctx context.Context, opts *ReturnListOptions) ([]Return, ReturnSummary, Pagination, error) {
	q := url.Values{}
	if opts != nil {
		opts.ListOptions.apply(q)
		if opts.Status != "" {
			q.Set("status", opts.Status)
		}
		if opts.Carrier != "" {
			q.Set("carrier", opts.Carrier)
		}
	}
	var out struct {
		listEnvelope[Return]
		Summary ReturnSummary `json:"summary"`
	}
	if err := s.client.get(ctx, "/api/v1/returns", q, &out); err != nil {
		return nil, ReturnSummary{}, Pagination{}, err
	}
	return out.Data, out.Summary, out.Pagination, nil
}
