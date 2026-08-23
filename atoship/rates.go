package atoship

import "context"

// RatesService quotes shipments.
type RatesService struct{ client *Client }

// RateRequest is one shipment to price.
type RateRequest struct {
	FromAddress AddressRef `json:"from_address"`
	ToAddress   AddressRef `json:"to_address"`
	Parcel      ParcelRef  `json:"parcel"`

	// CarrierAccountIDs restricts the quote to specific carrier accounts —
	// typically your own carrier contracts. Empty means every account available
	// to you.
	CarrierAccountIDs []string `json:"carrier_account_ids,omitempty"`
	// Carriers restricts by carrier, e.g. []string{"USPS", "UPS"}.
	Carriers []string `json:"carriers,omitempty"`
	// Services restricts by service name.
	Services []string `json:"services,omitempty"`

	Options *ShipmentOptions `json:"options,omitempty"`
}

type rateListResponse struct {
	listEnvelope[Rate]
	// AddressWarning is set when the destination looks deliverable but
	// imperfect — a corrected ZIP, a missing unit number. The rates are still
	// valid; the label may still be surcharged for the correction.
	AddressWarning any `json:"address_warning,omitempty"`
}

// List returns every available rate for one shipment, cheapest first.
//
// A rate ID is short-lived. Quote and buy in the same session rather than
// storing rates and buying later.
func (s *RatesService) List(ctx context.Context, req *RateRequest) ([]Rate, error) {
	var out rateListResponse
	if err := s.client.post(ctx, "/api/v1/rates", req, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// BatchRateRequest is one shipment inside a batch quote. ID is yours; it is
// echoed back so you can match a result to the row it came from.
type BatchRateRequest struct {
	ID string `json:"id"`
	RateRequest
}

// BatchRateResult is one shipment's outcome. Exactly one of Rates or Err is set.
type BatchRateResult struct {
	ID    string `json:"id"`
	Rates []Rate `json:"rates,omitempty"`
	Err   *Error `json:"error,omitempty"`
}

type batchRateResponse struct {
	Object    string            `json:"object"`
	Mode      string            `json:"mode"`
	Data      []BatchRateResult `json:"data"`
	CreatedAt string            `json:"created_at"`
}

// MaxBatchRateShipments is the server's per-request cap. Anything past it is
// dropped without an error, so split larger sets yourself.
const MaxBatchRateShipments = 50

// ListBatch quotes many shipments in one request.
//
// Use this instead of calling List in a loop. The per-organization rate limit is
// 5 requests per minute, so a page quoting twenty rows one at a time is blocked
// after the fifth; a batch is charged once no matter how many shipments it
// carries.
//
// One shipment failing does not fail the batch — that result carries Err and the
// rest still return rates. Results come back in the order submitted.
func (s *RatesService) ListBatch(ctx context.Context, shipments []BatchRateRequest) ([]BatchRateResult, error) {
	var out batchRateResponse
	body := struct {
		Shipments []BatchRateRequest `json:"shipments"`
	}{Shipments: shipments}
	if err := s.client.post(ctx, "/api/v1/rates/batch", body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
