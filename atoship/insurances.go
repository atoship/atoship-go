package atoship

import (
	"context"
	"net/url"
)

// InsurancesService insures shipments and files claims.
type InsurancesService struct{ client *Client }

// Insurance is coverage bought on a label.
type Insurance struct {
	ID             string  `json:"id"`
	Object         string  `json:"object"`
	LabelID        string  `json:"label_id"`
	TrackingNumber string  `json:"tracking_number"`
	Carrier        string  `json:"carrier"`
	InsuranceValue float64 `json:"insurance_value"`
	Premium        float64 `json:"premium"`
	Status         string  `json:"status"`
	Coverage       any     `json:"coverage"`
	CreatedAt      string  `json:"created_at"`
}

// InsuranceRequest insures an existing label for its declared value.
type InsuranceRequest struct {
	LabelID        string  `json:"label_id"`
	InsuranceValue float64 `json:"insurance_value"`
}

// Create insures a label. Premium is charged to your account.
func (s *InsurancesService) Create(ctx context.Context, req *InsuranceRequest) (*Insurance, error) {
	var out Insurance
	if err := s.client.post(ctx, "/api/v1/insurances", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one insurance policy.
func (s *InsurancesService) Get(ctx context.Context, id string) (*Insurance, error) {
	var out Insurance
	if err := s.client.get(ctx, "/api/v1/insurances/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List returns your insurance policies.
func (s *InsurancesService) List(ctx context.Context, opts *ListOptions) ([]Insurance, Pagination, error) {
	q := url.Values{}
	opts.apply(q)
	var out listEnvelope[Insurance]
	if err := s.client.get(ctx, "/api/v1/insurances", q, &out); err != nil {
		return nil, Pagination{}, err
	}
	return out.Data, out.Pagination, nil
}

// ClaimContact is who the insurer should reach about a claim.
type ClaimContact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// ClaimRequest files a claim against a policy.
//
// ClaimType is what went wrong — typically "damage" or "loss". IncidentDate
// matters: insurers enforce a window from the incident, not from when you
// noticed, so file with the real date rather than today's.
type ClaimRequest struct {
	ClaimType     string       `json:"claim_type"`
	ClaimAmount   float64      `json:"claim_amount"`
	Description   string       `json:"description"`
	IncidentDate  string       `json:"incident_date,omitempty"`
	EvidenceFiles []string     `json:"evidence_files,omitempty"`
	Contact       ClaimContact `json:"contact"`
}

// Claim is a filed insurance claim.
type Claim struct {
	ID             string   `json:"id"`
	Object         string   `json:"object"`
	LabelID        string   `json:"label_id"`
	TrackingNumber string   `json:"tracking_number"`
	ClaimSource    string   `json:"claim_source"`
	ClaimType      string   `json:"claim_type"`
	ClaimAmount    float64  `json:"claim_amount"`
	Status         string   `json:"status"`
	NextSteps      []string `json:"next_steps"`
	CreatedAt      string   `json:"created_at"`
}

// FileClaim files a claim against an insurance policy.
func (s *InsurancesService) FileClaim(ctx context.Context, insuranceID string, req *ClaimRequest) (*Claim, error) {
	var out Claim
	if err := s.client.post(ctx, "/api/v1/insurances/"+url.PathEscape(insuranceID)+"/claim", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListClaims returns the claims filed against a policy.
func (s *InsurancesService) ListClaims(ctx context.Context, insuranceID string) ([]Claim, error) {
	var out listEnvelope[Claim]
	if err := s.client.get(ctx, "/api/v1/insurances/"+url.PathEscape(insuranceID)+"/claim", nil, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// CancelClaim withdraws a claim you filed.
func (s *InsurancesService) CancelClaim(ctx context.Context, insuranceID, claimID string) error {
	path := "/api/v1/insurances/" + url.PathEscape(insuranceID) + "/claim/" + url.PathEscape(claimID) + "/cancel"
	return s.client.post(ctx, path, nil, nil)
}
