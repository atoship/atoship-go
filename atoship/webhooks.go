package atoship

import (
	"context"
	"net/url"
)

// WebhooksService manages the endpoints atoship posts events to.
type WebhooksService struct{ client *Client }

// Webhook is a registered endpoint.
type Webhook struct {
	ID        string   `json:"id"`
	Object    string   `json:"object"`
	URL       string   `json:"url"`
	Events    []string `json:"events"`
	IsActive  bool     `json:"is_active"`
	Secret    string   `json:"secret,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// WebhookRequest registers or updates an endpoint.
//
// Events empty means every event. Secret is returned once, on create — store it
// then; it is what you verify incoming deliveries against.
type WebhookRequest struct {
	URL      string   `json:"url,omitempty"`
	Events   []string `json:"events,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// List returns your registered webhooks.
func (s *WebhooksService) List(ctx context.Context, opts *ListOptions) ([]Webhook, Pagination, error) {
	q := url.Values{}
	opts.apply(q)
	var out listEnvelope[Webhook]
	if err := s.client.get(ctx, "/api/v1/webhooks", q, &out); err != nil {
		return nil, Pagination{}, err
	}
	return out.Data, out.Pagination, nil
}

// Get returns one webhook.
func (s *WebhooksService) Get(ctx context.Context, id string) (*Webhook, error) {
	var out Webhook
	if err := s.client.get(ctx, "/api/v1/webhooks/"+url.PathEscape(id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create registers an endpoint. The response carries Secret — this is the only
// time it is returned.
func (s *WebhooksService) Create(ctx context.Context, req *WebhookRequest) (*Webhook, error) {
	var out Webhook
	if err := s.client.post(ctx, "/api/v1/webhooks", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update replaces a webhook's configuration.
func (s *WebhooksService) Update(ctx context.Context, id string, req *WebhookRequest) (*Webhook, error) {
	var out Webhook
	if err := s.client.put(ctx, "/api/v1/webhooks/"+url.PathEscape(id), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Patch changes part of a webhook, leaving the rest alone.
func (s *WebhooksService) Patch(ctx context.Context, id string, req *WebhookRequest) (*Webhook, error) {
	var out Webhook
	if err := s.client.patch(ctx, "/api/v1/webhooks/"+url.PathEscape(id), req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete unregisters an endpoint.
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, "/api/v1/webhooks/"+url.PathEscape(id), nil)
}
