package atoship

import (
	"context"
	"time"
)

// TrackingService handles tracking-related operations
type TrackingService struct {
	client *Client
}

// TrackingInfo represents tracking information for a package
type TrackingInfo struct {
	TrackingNumber  string           `json:"trackingNumber"`
	Carrier         string           `json:"carrier"`
	Status          string           `json:"status"`
	EstimatedDelivery *time.Time     `json:"estimatedDelivery,omitempty"`
	ActualDelivery  *time.Time       `json:"actualDelivery,omitempty"`
	Events          []TrackingEvent  `json:"events"`
	CurrentLocation string           `json:"currentLocation,omitempty"`
	Delivered       bool             `json:"delivered"`
	Exception       bool             `json:"exception"`
	ExceptionReason string           `json:"exceptionReason,omitempty"`
}

// TrackingEvent represents a tracking event
type TrackingEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Details     string    `json:"details,omitempty"`
}

// Track tracks a package by tracking number
func (s *TrackingService) Track(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
	var info TrackingInfo
	err := s.client.get(ctx, "/api/tracking/"+trackingNumber, &info)
	return &info, err
}

// TrackWithCarrier tracks a package with a specific carrier
func (s *TrackingService) TrackWithCarrier(ctx context.Context, trackingNumber, carrier string) (*TrackingInfo, error) {
	var info TrackingInfo
	err := s.client.get(ctx, "/api/tracking/"+trackingNumber+"?carrier="+carrier, &info)
	return &info, err
}

// BatchTrack tracks multiple packages
func (s *TrackingService) BatchTrack(ctx context.Context, trackingNumbers []string) ([]TrackingInfo, error) {
	var infos []TrackingInfo
	req := map[string][]string{
		"trackingNumbers": trackingNumbers,
	}
	err := s.client.post(ctx, "/api/tracking/batch", req, &infos)
	return infos, err
}