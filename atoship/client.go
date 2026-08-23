// Package atoship is the Go client for the atoship shipping API.
//
// Quote a shipment, buy the cheapest label, print it:
//
//	client := atoship.NewClient(os.Getenv("ATOSHIP_API_KEY"))
//
//	rates, err := client.Rates.List(ctx, &atoship.RateRequest{
//		FromAddress: atoship.AddressRef{Inline: &atoship.Address{
//			Street1: "417 Montgomery St", City: "San Francisco", State: "CA", Zip: "94104",
//		}},
//		ToAddress: atoship.AddressRef{Inline: &atoship.Address{
//			Street1: "1600 Pennsylvania Ave NW", City: "Washington", State: "DC", Zip: "20500",
//		}},
//		Parcel: atoship.Parcel{Weight: 16, WeightUnit: "oz"},
//	})
//
//	label, err := client.Labels.Create(ctx, &atoship.LabelRequest{RateID: rates[0].ID})
//	label, err = client.Labels.Purchase(ctx, label.ID, nil)
//	fmt.Println(label.TrackingNumber, label.LabelURL)
//
// Every call takes a context and returns a typed value or an *Error carrying the
// API's own code, so you can branch on the reason:
//
//	if e, ok := err.(*atoship.Error); ok && e.Code == atoship.ErrInsufficientBalance {
//		…
//	}
//
// The client talks to /api/v1 only — the surface an API key can reach. It has no
// dependencies outside the standard library.
package atoship

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is where the API lives.
	DefaultBaseURL = "https://api.atoship.com"

	// DefaultTimeout applies per request. Buying a batch of labels is the slow
	// path — the server may spend a minute or more on 25 shipments — so raise it
	// with WithTimeout if you use Labels.BuyBatch.
	DefaultTimeout = 60 * time.Second

	// Version is this SDK's version, sent in the User-Agent.
	Version = "0.1.0"
)

// Client is an atoship API client. Safe for concurrent use.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	userAgent  string

	Account         *AccountService
	Addresses       *AddressesService
	CarrierAccounts *CarrierAccountsService
	Insurances      *InsurancesService
	Labels          *LabelsService
	Orders          *OrdersService
	Pickups         *PickupsService
	Rates           *RatesService
	Returns         *ReturnsService
	Tracking        *TrackingService
	Webhooks        *WebhooksService
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// NewClient returns a client authenticating with the given API key.
//
// A key beginning ak_test_ is a sandbox key: requests are validated and priced
// but nothing is bought and nothing is charged.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: DefaultTimeout},
		userAgent:  "atoship-go/" + Version,
	}
	for _, opt := range opts {
		opt(c)
	}

	c.Account = &AccountService{client: c}
	c.Addresses = &AddressesService{client: c}
	c.CarrierAccounts = &CarrierAccountsService{client: c}
	c.Insurances = &InsurancesService{client: c}
	c.Labels = &LabelsService{client: c}
	c.Orders = &OrdersService{client: c}
	c.Pickups = &PickupsService{client: c}
	c.Rates = &RatesService{client: c}
	c.Returns = &ReturnsService{client: c}
	c.Tracking = &TrackingService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	return c
}

// WithBaseURL points the client at a different host.
func WithBaseURL(raw string) ClientOption {
	return func(c *Client) { c.baseURL = strings.TrimRight(raw, "/") }
}

// WithTimeout sets the per-request timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// WithHTTPClient supplies your own *http.Client — for a custom transport,
// proxy, or instrumentation. Its Timeout is used as-is, so set one.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithUserAgent appends your application's name to the User-Agent, which is
// what support sees when tracing a request.
func WithUserAgent(product string) ClientOption {
	return func(c *Client) {
		if product != "" {
			c.userAgent = product + " " + c.userAgent
		}
	}
}

// Error is a failure reported by the API. Compare Code against the Err*
// constants rather than matching on Message, which is prose and may change.
type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	// Details carries whatever extra the endpoint attached, e.g. per-field
	// validation failures.
	Details any `json:"details,omitempty"`
}

func (e *Error) Error() string {
	if e.Code == "" {
		return fmt.Sprintf("atoship: %s (HTTP %d)", e.Message, e.StatusCode)
	}
	return fmt.Sprintf("atoship: %s [%s] (HTTP %d)", e.Message, e.Code, e.StatusCode)
}

// Error codes worth branching on.
const (
	ErrUnauthorized        = "UNAUTHORIZED"
	ErrValidation          = "VALIDATION_ERROR"
	ErrNotFound            = "NOT_FOUND"
	ErrRateLimited         = "RATE_LIMIT_EXCEEDED"
	ErrInsufficientBalance = "INSUFFICIENT_BALANCE"
	ErrPaymentRequired     = "PAYMENT_REQUIRED"
	ErrInternal            = "INTERNAL_ERROR"
)

// IsNotFound reports whether err is an API 404.
func IsNotFound(err error) bool { return codeIs(err, ErrNotFound) }

// IsRateLimited reports whether err is the per-organization rate limit.
//
// That limit is 5 requests per minute. If you are hitting it while quoting or
// buying many shipments, the fix is not to back off — it is Rates.ListBatch and
// Labels.BuyBatch, which cost one request for the whole set.
func IsRateLimited(err error) bool { return codeIs(err, ErrRateLimited) }

// IsUnauthorized reports whether the API key was missing, malformed or revoked.
func IsUnauthorized(err error) bool { return codeIs(err, ErrUnauthorized) }

func codeIs(err error, code string) bool {
	e, ok := err.(*Error)
	return ok && e.Code == code
}

// errorEnvelope is how every v1 route reports a failure. The `object` field is
// spelled both "Error" and "error" by different routes, so it is not used as the
// discriminator — the presence of a code is.
type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details any    `json:"details,omitempty"`
	} `json:"error"`
}

// SandboxError is returned when a sandbox key (ak_test_…) is used for an
// operation that would spend money. The request was validated and rejected
// before execution — nothing was bought and nothing was charged.
//
// This is its own type because the API answers a sandbox write with a different
// body entirely:
//
//	{"success":true,"sandbox":true,"message":"…","validatedFields":[…]}
//
// no object, no id, no status. Decoded into the normal result type that is a
// zero-valued struct with a nil error — a Label with no ID that looks like it
// was created. Failing loudly here is the whole point: a caller who sees this
// knows their request shape is correct and that only the key needs swapping.
type SandboxError struct {
	Message         string   `json:"message"`
	ValidatedFields []string `json:"validatedFields"`
	Hints           []string `json:"hints"`
}

func (e *SandboxError) Error() string {
	return "atoship: sandbox key — request validated but not executed: " + e.Message
}

// IsSandbox reports whether err is a sandbox key refusing to spend money.
// The request itself was well-formed.
func IsSandbox(err error) bool {
	_, ok := err.(*SandboxError)
	return ok
}

// sandboxEnvelope is how a sandbox key answers a write.
//
// `sandbox: true` alone is NOT the marker. /v1/labels/batch answers a sandbox
// request with a full, useful manifest — {"object":"BatchResult","mode":
// "sandbox","sandbox":true,summary,results} — and an earlier version of this
// check swallowed that into an error, throwing away a per-shipment dry run the
// caller wanted. The distinguishing fact is the absence of `object`: a body
// with no object carries no result to decode, so a zero-valued struct is all
// the caller would otherwise get.
type sandboxEnvelope struct {
	Success         bool     `json:"success"`
	Sandbox         bool     `json:"sandbox"`
	Object          string   `json:"object"`
	Message         string   `json:"message"`
	ValidatedFields []string `json:"validatedFields"`
	Hints           []string `json:"hints"`
}

// Pagination accompanies every list response.
type Pagination struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	Pages   int  `json:"pages"`
	HasMore bool `json:"has_more"`
}

// ListOptions are the query parameters every list endpoint accepts.
type ListOptions struct {
	Page  int
	Limit int
}

func (o *ListOptions) apply(q url.Values) {
	if o == nil {
		return
	}
	if o.Page > 0 {
		q.Set("page", fmt.Sprint(o.Page))
	}
	if o.Limit > 0 {
		q.Set("limit", fmt.Sprint(o.Limit))
	}
}

// do performs one request and decodes the whole response body into out.
//
// The body is NOT unwrapped from a `data` envelope here: some endpoints return
// the object itself (`{"object":"Label","id":…}`) and some return a list
// (`{"object":"LabelList","data":[…]}`), so each caller passes a type matching
// the shape it expects. An earlier version of this SDK assumed a
// `{"success":true,"data":…}` envelope that this API has never used, which
// turned every successful call into an error.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("atoship: encoding request body: %w", err)
		}
		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return fmt.Errorf("atoship: building request: %w", err)
	}
	// The API reads the key from Authorization only. X-API-Key is ignored, which
	// is what the previous version of this SDK sent.
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("atoship: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("atoship: reading response: %w", err)
	}

	// An error can arrive with a 2xx in one place — batch rows carry their own —
	// but at the envelope level a failure always has both a non-2xx status and a
	// code. Check the status first so an empty body still produces an error.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &Error{StatusCode: resp.StatusCode}
		var env errorEnvelope
		if json.Unmarshal(raw, &env) == nil && env.Error.Code != "" {
			apiErr.Code = env.Error.Code
			apiErr.Message = env.Error.Message
			apiErr.Details = env.Error.Details
		} else {
			apiErr.Message = strings.TrimSpace(string(raw))
			if apiErr.Message == "" {
				apiErr.Message = http.StatusText(resp.StatusCode)
			}
		}
		return apiErr
	}

	// A sandbox key answers a money-spending write with its own envelope. Catch
	// it before decoding, or it unmarshals into a zero-valued result and the
	// caller gets an object that never existed with no error to tell them.
	var sb sandboxEnvelope
	if json.Unmarshal(raw, &sb) == nil && sb.Sandbox && sb.Object == "" {
		return &SandboxError{Message: sb.Message, ValidatedFields: sb.ValidatedFields, Hints: sb.Hints}
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("atoship: decoding %s %s response: %w", method, path, err)
	}
	return nil
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, out)
}

func (c *Client) post(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPost, path, nil, body, out)
}

func (c *Client) put(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPut, path, nil, body, out)
}

func (c *Client) patch(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPatch, path, nil, body, out)
}

func (c *Client) del(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil, out)
}
