package atoship

import (
	"context"
	"time"
)

// ShippingService handles shipping-related operations
type ShippingService struct {
	client *Client
}

// Parcel represents a package
type Parcel struct {
	Length     float64 `json:"length"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	DimUnit    string  `json:"dimUnit"`
	Weight     float64 `json:"weight"`
	WeightUnit string  `json:"weightUnit"`
}

// RateRequest represents a request for shipping rates
type RateRequest struct {
	FromAddress *Address `json:"fromAddress"`
	ToAddress   *Address `json:"toAddress"`
	Parcel      *Parcel  `json:"parcel"`
	ShipDate    string   `json:"shipDate,omitempty"`
	Insurance   float64  `json:"insurance,omitempty"`
}

// ShippingRate represents a shipping rate
type ShippingRate struct {
	ID             string    `json:"id"`
	Carrier        string    `json:"carrier"`
	Service        string    `json:"service"`
	ServiceCode    string    `json:"serviceCode"`
	Rate           float64   `json:"rate"`
	Currency       string    `json:"currency"`
	DeliveryDays   int       `json:"deliveryDays,omitempty"`
	DeliveryDate   time.Time `json:"deliveryDate,omitempty"`
	Insurance      float64   `json:"insurance,omitempty"`
	Tracking       bool      `json:"tracking"`
}

// PurchaseLabelRequest represents a request to purchase a shipping label
type PurchaseLabelRequest struct {
	RateID        string         `json:"rateId"`
	OrderID       string         `json:"orderId,omitempty"`
	LabelFormat   string         `json:"labelFormat,omitempty"`
	Customs       *CustomsInfo   `json:"customs,omitempty"`
	ReturnLabel   bool           `json:"returnLabel,omitempty"`
}

// ShippingLabel represents a shipping label
type ShippingLabel struct {
	ID              string    `json:"id"`
	TrackingNumber  string    `json:"trackingNumber"`
	LabelURL        string    `json:"labelUrl"`
	LabelPDF        string    `json:"labelPdf,omitempty"`
	Carrier         string    `json:"carrier"`
	Service         string    `json:"service"`
	Rate            float64   `json:"rate"`
	CreatedAt       time.Time `json:"createdAt"`
}

// CustomsInfo represents customs information
type CustomsInfo struct {
	ContentsType    string         `json:"contentsType"`
	ContentsExplanation string     `json:"contentsExplanation,omitempty"`
	CustomsItems    []CustomsItem  `json:"customsItems"`
}

// CustomsItem represents a customs item
type CustomsItem struct {
	Description     string  `json:"description"`
	Quantity        int     `json:"quantity"`
	Value           float64 `json:"value"`
	Weight          float64 `json:"weight"`
	OriginCountry   string  `json:"originCountry"`
	HSTariffNumber  string  `json:"hsTariffNumber,omitempty"`
}

// GetRates gets shipping rates for a package
func (s *ShippingService) GetRates(ctx context.Context, req *RateRequest) ([]ShippingRate, error) {
	var rates []ShippingRate
	err := s.client.post(ctx, "/api/carriers/smart-rates", req, &rates)
	return rates, err
}

// PurchaseLabel purchases a shipping label using V2 API with routing engine
func (s *ShippingService) PurchaseLabel(ctx context.Context, req *PurchaseLabelRequest) (*ShippingLabel, error) {
	var label ShippingLabel
	err := s.client.post(ctx, "/api/labels/purchase-v2", req, &label)
	return &label, err
}

// GetLabel retrieves a shipping label by ID
func (s *ShippingService) GetLabel(ctx context.Context, labelID string) (*ShippingLabel, error) {
	var label ShippingLabel
	err := s.client.get(ctx, "/api/labels/"+labelID, &label)
	return &label, err
}

// CancelLabel cancels a shipping label
func (s *ShippingService) CancelLabel(ctx context.Context, labelID string) (*ShippingLabel, error) {
	var label ShippingLabel
	err := s.client.post(ctx, "/api/labels/"+labelID+"/cancel", nil, &label)
	return &label, err
}