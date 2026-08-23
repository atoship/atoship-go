package atoship

import (
	"context"
	"fmt"
	"net/url"
)

// LabelsService creates, buys, lists and voids shipping labels.
type LabelsService struct{ client *Client }

// LabelRequest creates a label. Nothing is charged until Purchase.
//
// Either give RateID from a quote, or name Carrier and Service and let the
// server price it. RateID is the safer path: it buys the exact price you showed
// the customer.
type LabelRequest struct {
	RateID string `json:"rate_id,omitempty"`

	FromAddress AddressRef `json:"from_address"`
	ToAddress   AddressRef `json:"to_address"`
	Parcel      ParcelRef  `json:"parcel"`

	Carrier          string `json:"carrier,omitempty"`
	Service          string `json:"service,omitempty"`
	ServiceCode      string `json:"service_code,omitempty"`
	CarrierAccountID string `json:"carrier_account_id,omitempty"`

	LabelFormat string `json:"label_format,omitempty"` // pdf (default), png, zpl
	LabelSize   string `json:"label_size,omitempty"`

	// Reference is your own identifier, echoed back on the label and in lists.
	Reference string `json:"reference,omitempty"`

	// Order linkage, when the shipment came from a connected store.
	Platform        string `json:"platform,omitempty"`
	PlatformOrderID string `json:"platform_order_id,omitempty"`
	Channel         string `json:"channel,omitempty"`
	ExternalOrderID string `json:"external_order_id,omitempty"`
	OrderID         string `json:"order_id,omitempty"`

	LabelCustomFields *LabelCustomFields `json:"label_custom_fields,omitempty"`
	Options           *ShipmentOptions   `json:"options,omitempty"`
}

// Create makes a label without buying it. The returned Label has an ID to pass
// to Purchase.
func (s *LabelsService) Create(ctx context.Context, req *LabelRequest) (*Label, error) {
	var out Label
	if err := s.client.post(ctx, "/api/v1/labels", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PurchaseRequest buys a created label. Customs is required for an
// international shipment and ignored for a domestic one.
type PurchaseRequest struct {
	Customs *Customs `json:"customs,omitempty"`
}

// Customs declares the contents of an international shipment.
//
// Set Incoterm to "DAP" unless you have decided otherwise: without it the duty
// is billed to the shipper — us — and that only surfaces weeks later on the
// carrier invoice.
type Customs struct {
	ContentsType        string        `json:"contents_type,omitempty"`
	ContentsExplanation string        `json:"contents_explanation,omitempty"`
	Incoterm            string        `json:"incoterm,omitempty"`
	NonDelivery         string        `json:"non_delivery_option,omitempty"`
	Certify             bool          `json:"certify,omitempty"`
	CertifySigner       string        `json:"certify_signer,omitempty"`
	Items               []CustomsItem `json:"items,omitempty"`
}

// CustomsItem is one line on the customs declaration.
type CustomsItem struct {
	Description    string  `json:"description"`
	Quantity       int     `json:"quantity"`
	Value          float64 `json:"value"`
	Weight         float64 `json:"weight,omitempty"`
	WeightUnit     string  `json:"weight_unit,omitempty"`
	HSTariffNumber string  `json:"hs_tariff_number,omitempty"`
	OriginCountry  string  `json:"origin_country,omitempty"`
	Currency       string  `json:"currency,omitempty"`
}

// Purchase buys a created label. This is the call that charges your account.
//
// The returned Label may carry an empty TrackingNumber with TrackingPending
// set: on some services the carrier assigns the number asynchronously. That is
// not a failure — the label is bought and printable. Poll Get until the number
// appears.
func (s *LabelsService) Purchase(ctx context.Context, labelID string, req *PurchaseRequest) (*Label, error) {
	if req == nil {
		req = &PurchaseRequest{}
	}
	var out Label
	if err := s.client.post(ctx, "/api/v1/labels/"+url.PathEscape(labelID)+"/purchase", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one label.
func (s *LabelsService) Get(ctx context.Context, labelID string) (*Label, error) {
	var out Label
	if err := s.client.get(ctx, "/api/v1/labels/"+url.PathEscape(labelID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Void cancels a label and refunds it where the carrier allows. Carriers have
// their own windows — a void accepted by the API can still be refused by the
// carrier days later.
func (s *LabelsService) Void(ctx context.Context, labelID string) error {
	return s.client.del(ctx, "/api/v1/labels/"+url.PathEscape(labelID), nil)
}

// LabelListOptions filters a label list.
type LabelListOptions struct {
	ListOptions
	Status    string
	Carrier   string
	Reference string
}

// List returns labels newest first, with the pagination block.
func (s *LabelsService) List(ctx context.Context, opts *LabelListOptions) ([]Label, Pagination, error) {
	q := url.Values{}
	if opts != nil {
		opts.ListOptions.apply(q)
		if opts.Status != "" {
			q.Set("status", opts.Status)
		}
		if opts.Carrier != "" {
			q.Set("carrier", opts.Carrier)
		}
		if opts.Reference != "" {
			q.Set("reference", opts.Reference)
		}
	}
	var out listEnvelope[Label]
	if err := s.client.get(ctx, "/api/v1/labels", q, &out); err != nil {
		return nil, Pagination{}, err
	}
	return out.Data, out.Pagination, nil
}

// DownloadURL is the atoship-hosted PDF for a label. It is also returned as
// Label.LabelURL; this builds it without another API call.
//
// The URL carries a per-label token, so it can be opened in a browser or
// embedded directly.
func (s *LabelsService) DownloadURL(labelID string) string {
	return s.client.baseURL + "/api/v1/labels/" + url.PathEscape(labelID) + "/download"
}

// MaxBatchLabelShipments is the server's per-request cap.
const MaxBatchLabelShipments = 25

// BatchShipment is one shipment in a buy-batch.
//
// Reference must be unique within the batch. It identifies the row in the
// result and is what makes a resubmit safe.
type BatchShipment struct {
	Reference   string     `json:"reference"`
	FromAddress AddressRef `json:"from_address,omitempty"`
	ToAddress   AddressRef `json:"to_address"`
	Parcel      ParcelRef  `json:"parcel"`
}

// BatchServiceSelection names one carrier and service for a whole batch.
type BatchServiceSelection struct {
	Carrier string `json:"carrier"`
	Service string `json:"service"`
}

// BuyBatchRequest buys labels for many shipments in one request.
//
// Domestic only: this path carries no customs data, so an international
// shipment is rejected before anything is bought. Insurance, signature and
// hazmat are likewise single-label features. Anything outside that goes through
// Create + Purchase, which is fully featured.
type BuyBatchRequest struct {
	// Select is "cheapest" (the default) or "service".
	Select string `json:"select,omitempty"`
	// Service is required when Select is "service", and applies to every
	// shipment in the batch. There is no per-shipment service here.
	Service   *BatchServiceSelection `json:"service,omitempty"`
	Shipments []BatchShipment        `json:"shipments"`
}

// BatchRow is one shipment's outcome.
//
// Status "needs_review" means the wallet was charged but the label was not
// recorded. That row is held for reconciliation — do not re-buy it.
type BatchRow struct {
	Index          int    `json:"index"`
	Reference      string `json:"reference"`
	Status         string `json:"status"`
	Carrier        string `json:"carrier"`
	Service        string `json:"service"`
	TrackingNumber string `json:"tracking_number"`
	// Cost is the amount charged to your account for this label.
	Cost      float64 `json:"cost"`
	LabelID   string  `json:"label_id"`
	LabelURL  string  `json:"label_url"`
	ErrorCode string  `json:"error_code"`
	Error     string  `json:"error"`
}

// BatchSummary counts the rows by outcome.
type BatchSummary struct {
	Total        int     `json:"total"`
	Purchased    int     `json:"purchased"`
	Failed       int     `json:"failed"`
	Skipped      int     `json:"skipped"`
	NeedsReview  int     `json:"needs_review"`
	PriceChanged int     `json:"price_changed"`
	Error        int     `json:"error"`
	ChargedTotal float64 `json:"charged_total"`
}

// BatchResult is a batch manifest.
type BatchResult struct {
	Object  string       `json:"object"`
	Mode    string       `json:"mode"`
	BatchID string       `json:"batch_id"`
	Running bool         `json:"running"`
	Summary BatchSummary `json:"summary"`
	Results []BatchRow   `json:"results"`
}

// BuyBatch buys up to MaxBatchLabelShipments labels in one request.
//
// Bought one at a time, each shipment costs three requests — quote, create,
// purchase — and the per-organization limit of 5 requests per minute stops a
// warehouse app after the second parcel. A batch is one request for the whole
// selection.
//
// A nil error does not mean every shipment succeeded. Read Summary and each
// row's Status.
//
// If this call times out or the connection drops, do NOT resubmit blind — call
// GetBatch with the BatchID. The stored manifest is authoritative and a resubmit
// can buy a second label.
func (s *LabelsService) BuyBatch(ctx context.Context, req *BuyBatchRequest) (*BatchResult, error) {
	if len(req.Shipments) > MaxBatchLabelShipments {
		return nil, fmt.Errorf("atoship: %d shipments exceeds the batch limit of %d", len(req.Shipments), MaxBatchLabelShipments)
	}
	var out BatchResult
	if err := s.client.post(ctx, "/api/v1/labels/batch", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBatch re-reads a batch manifest. Use it after a timeout to find out what
// was actually bought, instead of resubmitting.
func (s *LabelsService) GetBatch(ctx context.Context, batchID string) (*BatchResult, error) {
	q := url.Values{"batch_id": {batchID}}
	var out BatchResult
	if err := s.client.get(ctx, "/api/v1/labels/batch", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
