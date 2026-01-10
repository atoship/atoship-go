package atoship

import (
	"context"
	"fmt"
)

// WebhooksService handles webhook-related operations
type WebhooksService struct {
	client *Client
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID      string   `json:"id"`
	URL     string   `json:"url"`
	Events  []string `json:"events"`
	Active  bool     `json:"active"`
	Secret  string   `json:"secret,omitempty"`
}

// CreateWebhookRequest represents a request to create a webhook
type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Active bool     `json:"active,omitempty"`
}

// Create creates a new webhook
func (s *WebhooksService) Create(ctx context.Context, req *CreateWebhookRequest) (*Webhook, error) {
	var webhook Webhook
	err := s.client.post(ctx, "/api/admin/webhooks", req, &webhook)
	return &webhook, err
}

// List lists all webhooks
func (s *WebhooksService) List(ctx context.Context) ([]Webhook, error) {
	var webhooks []Webhook
	err := s.client.get(ctx, "/api/admin/webhooks", &webhooks)
	return webhooks, err
}

// Delete deletes a webhook
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error {
	return s.client.delete(ctx, fmt.Sprintf("/api/admin/webhooks/%s", webhookID))
}