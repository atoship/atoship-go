package atoship

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultBaseURL is the default base URL for the atoship API
	DefaultBaseURL = "https://api.atoship.com"
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
	// Version is the SDK version
	Version = "1.0.0"
)

// Client is the main atoship API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *resty.Client
	debug      bool

	// Services
	Orders    *OrdersService
	Addresses *AddressesService
	Shipping  *ShippingService
	Tracking  *TrackingService
	Users     *UsersService
	Admin     *AdminService
	Carriers  *CarriersService
	Webhooks  *WebhooksService
}

// ClientOption is a function that configures the client
type ClientOption func(*Client)

// NewClient creates a new atoship API client
func NewClient(apiKey string, opts ...ClientOption) *Client {
	client := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: resty.New().
			SetTimeout(DefaultTimeout).
			SetHeader("User-Agent", "atoship-go-sdk/"+Version),
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	// Set up authentication
	client.httpClient.SetHeader("X-API-Key", apiKey)
	client.httpClient.SetHeader("Content-Type", "application/json")

	// Set debug mode
	if client.debug {
		client.httpClient.SetDebug(true)
	}

	// Initialize services
	client.Orders = &OrdersService{client: client}
	client.Addresses = &AddressesService{client: client}
	client.Shipping = &ShippingService{client: client}
	client.Tracking = &TrackingService{client: client}
	client.Users = &UsersService{client: client}
	client.Admin = &AdminService{client: client}
	client.Carriers = &CarriersService{client: client}
	client.Webhooks = &WebhooksService{client: client}

	return client
}

// WithBaseURL sets a custom base URL
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
		c.httpClient.SetBaseURL(url)
	}
}

// WithTimeout sets a custom request timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.SetTimeout(timeout)
	}
}

// WithRetryCount sets the number of retry attempts
func WithRetryCount(count int) ClientOption {
	return func(c *Client) {
		c.httpClient.SetRetryCount(count)
	}
}

// WithDebug enables debug mode
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
		c.httpClient.SetDebug(debug)
	}
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool            `json:"success"`
	Data      json.RawMessage `json:"data,omitempty"`
	Error     string          `json:"error,omitempty"`
	RequestID string          `json:"requestId,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Items   json.RawMessage `json:"items"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
	HasMore bool            `json:"hasMore"`
}

// APIError represents an API error
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	RequestID  string `json:"requestId"`
	Details    any    `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("atoship API error: %s (code: %s, status: %d)", e.Message, e.Code, e.StatusCode)
}

// Common error codes
const (
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeAuthentication  = "AUTHENTICATION_ERROR"
	ErrCodeAuthorization   = "AUTHORIZATION_ERROR"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeRateLimit       = "RATE_LIMIT_EXCEEDED"
	ErrCodeServerError     = "SERVER_ERROR"
	ErrCodeNetworkError    = "NETWORK_ERROR"
	ErrCodeTimeoutError    = "TIMEOUT_ERROR"
	ErrCodeConfigError     = "CONFIGURATION_ERROR"
)

// makeRequest performs an HTTP request
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	req := c.httpClient.R().
		SetContext(ctx).
		SetResult(&APIResponse{})

	if body != nil {
		req.SetBody(body)
	}

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "DELETE":
		resp, err = req.Delete(path)
	case "PATCH":
		resp, err = req.Patch(path)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return &APIError{
			Code:    ErrCodeNetworkError,
			Message: err.Error(),
		}
	}

	// Check for HTTP errors
	if resp.IsError() {
		var apiErr APIError
		if err := json.Unmarshal(resp.Body(), &apiErr); err != nil {
			return &APIError{
				Code:       ErrCodeServerError,
				Message:    string(resp.Body()),
				StatusCode: resp.StatusCode(),
			}
		}
		apiErr.StatusCode = resp.StatusCode()
		return &apiErr
	}

	// Parse successful response
	apiResp := resp.Result().(*APIResponse)
	if !apiResp.Success {
		return &APIError{
			Code:      "API_ERROR",
			Message:   apiResp.Error,
			RequestID: apiResp.RequestID,
		}
	}

	// Unmarshal data into result
	if result != nil && apiResp.Data != nil {
		return json.Unmarshal(apiResp.Data, result)
	}

	return nil
}

// get performs a GET request
func (c *Client) get(ctx context.Context, path string, result interface{}) error {
	return c.makeRequest(ctx, "GET", path, nil, result)
}

// post performs a POST request
func (c *Client) post(ctx context.Context, path string, body, result interface{}) error {
	return c.makeRequest(ctx, "POST", path, body, result)
}

// put performs a PUT request
func (c *Client) put(ctx context.Context, path string, body, result interface{}) error {
	return c.makeRequest(ctx, "PUT", path, body, result)
}

// delete performs a DELETE request
func (c *Client) delete(ctx context.Context, path string) error {
	return c.makeRequest(ctx, "DELETE", path, nil, nil)
}

// patch performs a PATCH request
func (c *Client) patch(ctx context.Context, path string, body, result interface{}) error {
	return c.makeRequest(ctx, "PATCH", path, body, result)
}