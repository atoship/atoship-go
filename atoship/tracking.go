package atoship

import (
	"context"
	"net/url"
)

// TrackingService follows shipments.
type TrackingService struct{ client *Client }

// TrackingEvent is one scan.
type TrackingEvent struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Timestamp   string `json:"timestamp"`
}

// Tracking is a shipment's current state and its scan history.
//
// Status is normalised across carriers: pre_transit, in_transit, out_for_delivery,
// delivered, available_for_pickup, return_to_sender, failure, cancelled, unknown.
// StatusDetail is the carrier's own wording and varies.
type Tracking struct {
	Object            string          `json:"object"`
	TrackingNumber    string          `json:"tracking_number"`
	Carrier           string          `json:"carrier"`
	Service           string          `json:"service"`
	Status            string          `json:"status"`
	StatusDetail      string          `json:"status_detail"`
	EstimatedDelivery string          `json:"estimated_delivery"`
	ActualDelivery    string          `json:"actual_delivery"`
	Origin            Address         `json:"origin"`
	Destination       Address         `json:"destination"`
	Events            []TrackingEvent `json:"events"`
	UpdatedAt         string          `json:"updated_at"`
}

// Get returns tracking for a number.
//
// A label bought minutes ago may have no scans yet: the carrier has the number
// but has not received the parcel. That returns status pre_transit with an empty
// Events, which is normal and not an error.
func (s *TrackingService) Get(ctx context.Context, trackingNumber string) (*Tracking, error) {
	var out Tracking
	if err := s.client.get(ctx, "/api/v1/tracking/"+url.PathEscape(trackingNumber), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
