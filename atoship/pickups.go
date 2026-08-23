package atoship

import (
	"context"
	"net/url"
)

// PickupsService schedules carrier pickups.
type PickupsService struct{ client *Client }

// PickupTimeWindow is when the driver may come.
type PickupTimeWindow struct {
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// Pickup is a scheduled carrier collection.
type Pickup struct {
	ID                 string  `json:"id"`
	Object             string  `json:"object"`
	Status             string  `json:"status"`
	ConfirmationNumber string  `json:"confirmation_number"`
	Carrier            string  `json:"carrier"`
	PickupDate         string  `json:"pickup_date"`
	PickupAddress      Address `json:"pickup_address"`
	PackageCount       int     `json:"package_count"`
	TotalWeight        float64 `json:"total_weight"`
	// EstimatedCost is what the pickup will cost YOU. Scheduled pickups on a
	// regular route are free; an on-demand pickup carries a fee.
	EstimatedCost       float64 `json:"estimated_cost"`
	SpecialInstructions string  `json:"special_instructions"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
	CancelledAt         string  `json:"cancelled_at"`
}

// PickupRequest schedules a pickup.
//
// PickupDate is a date, not a timestamp. Carriers have same-day cutoffs and
// will reject a request made after theirs, so schedule the day before where you
// can.
type PickupRequest struct {
	PickupAddress       Address           `json:"pickup_address"`
	PickupDate          string            `json:"pickup_date"`
	PickupTimeWindow    *PickupTimeWindow `json:"pickup_time_window,omitempty"`
	Packages            []string          `json:"packages,omitempty"`
	PackageCount        int               `json:"package_count,omitempty"`
	TotalWeight         float64           `json:"total_weight,omitempty"`
	SpecialInstructions string            `json:"special_instructions,omitempty"`
	Carrier             string            `json:"carrier,omitempty"`
}

// Create schedules a pickup. The returned ConfirmationNumber is the carrier's —
// it is what their support desk asks for.
func (s *PickupsService) Create(ctx context.Context, req *PickupRequest) (*Pickup, error) {
	var out Pickup
	if err := s.client.post(ctx, "/api/v1/pickups", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one pickup.
func (s *PickupsService) Get(ctx context.Context, id string) (*Pickup, error) {
	var out Pickup
	if err := s.client.get(ctx, "/api/v1/pickups/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List returns your scheduled pickups.
func (s *PickupsService) List(ctx context.Context, opts *ListOptions) ([]Pickup, Pagination, error) {
	q := url.Values{}
	opts.apply(q)
	var out listEnvelope[Pickup]
	if err := s.client.get(ctx, "/api/v1/pickups", q, &out); err != nil {
		return nil, Pagination{}, err
	}
	return out.Data, out.Pagination, nil
}

// Cancel calls off a scheduled pickup. Carriers stop accepting cancellations
// close to the window, so a late cancel can be refused.
func (s *PickupsService) Cancel(ctx context.Context, id string) (*Pickup, error) {
	var out Pickup
	if err := s.client.post(ctx, "/api/v1/pickups/"+url.PathEscape(id)+"/cancel", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
