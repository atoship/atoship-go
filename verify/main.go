// Integration check: every SDK method against the live v1 API.
//
// TEMPORARY — not part of the published SDK. The point is that the previous
// version of this SDK was never run against the API even once, and three
// independent faults (wrong auth header, wrong endpoints, wrong response
// envelope) each made every call fail. Nothing here is stubbed.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atoship/atoship-go/atoship"
)

var (
	pass, fail, skip int
	ctx              context.Context
)

func check(name string, fn func() (string, error)) {
	detail, err := fn()
	switch {
	// SKIP is tested FIRST. With `err == nil` first, a skip returns a nil error
	// and is counted as a pass — which is how the first run of this harness
	// reported "17 passed, 0 skipped" while one endpoint had never been called.
	case strings.HasPrefix(detail, "SKIP"):
		skip++
		fmt.Printf("  \033[33mSKIP\033[0m  %-34s %s\n", name, strings.TrimPrefix(detail, "SKIP "))
	case err == nil:
		pass++
		fmt.Printf("  \033[32mPASS\033[0m  %-34s %s\n", name, detail)
	default:
		fail++
		fmt.Printf("  \033[31mFAIL\033[0m  %-34s %v\n", name, err)
	}
}

func main() {
	key := os.Getenv("ATOSHIP_API_KEY")
	if key == "" {
		fmt.Println("set ATOSHIP_API_KEY")
		os.Exit(2)
	}
	base := os.Getenv("ATOSHIP_BASE_URL")
	opts := []atoship.ClientOption{atoship.WithTimeout(90 * time.Second)}
	if base != "" {
		opts = append(opts, atoship.WithBaseURL(base))
	}
	c := atoship.NewClient(key, opts...)
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	from := atoship.InlineAddress(atoship.Address{
		Name: "atoship", Street1: "417 Montgomery St", City: "San Francisco",
		State: "CA", Zip: "94104", Country: "US",
	})
	to := atoship.InlineAddress(atoship.Address{
		Name: "Test Recipient", Street1: "1600 Pennsylvania Ave NW", City: "Washington",
		State: "DC", Zip: "20500", Country: "US",
	})
	parcel := atoship.InlineParcel(atoship.Parcel{
		Length: 10, Width: 8, Height: 4, DimensionUnit: "in", Weight: 16, WeightUnit: "oz",
	})

	fmt.Println("\nREAD")
	check("Account.Get", func() (string, error) {
		a, err := c.Account.Get(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s  $%.2f available", a.Email, a.Balance.Available), nil
	})
	check("Addresses.List", func() (string, error) {
		a, p, err := c.Addresses.List(ctx, &atoship.ListOptions{Limit: 3})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d of %d", len(a), p.Total), nil
	})
	check("CarrierAccounts.List", func() (string, error) {
		a, err := c.CarrierAccounts.List(ctx, "")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d accounts", len(a)), nil
	})
	check("Webhooks.List", func() (string, error) {
		w, p, err := c.Webhooks.List(ctx, &atoship.ListOptions{Limit: 3})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d of %d", len(w), p.Total), nil
	})
	check("Returns.List", func() (string, error) {
		r, s, _, err := c.Returns.List(ctx, &atoship.ReturnListOptions{ListOptions: atoship.ListOptions{Limit: 3}})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d rows, summary total_cost $%.2f", len(r), s.TotalCost), nil
	})
	check("Insurances.List", func() (string, error) {
		i, p, err := c.Insurances.List(ctx, &atoship.ListOptions{Limit: 3})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d of %d", len(i), p.Total), nil
	})
	check("Pickups.List", func() (string, error) {
		p, pg, err := c.Pickups.List(ctx, &atoship.ListOptions{Limit: 3})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d of %d", len(p), pg.Total), nil
	})

	var labelID, trackingNumber string
	check("Labels.List", func() (string, error) {
		l, p, err := c.Labels.List(ctx, &atoship.LabelListOptions{ListOptions: atoship.ListOptions{Limit: 3}})
		if err != nil {
			return "", err
		}
		for _, x := range l {
			if labelID == "" {
				labelID = x.ID
			}
			if x.TrackingNumber != "" && trackingNumber == "" {
				trackingNumber = x.TrackingNumber
			}
		}
		return fmt.Sprintf("%d of %d", len(l), p.Total), nil
	})
	check("Labels.Get", func() (string, error) {
		if labelID == "" {
			return "SKIP no label to read", nil
		}
		l, err := c.Labels.Get(ctx, labelID)
		if err != nil {
			return "", err
		}
		rate := "null"
		if l.Rate != nil {
			rate = fmt.Sprintf("$%.2f", *l.Rate)
		}
		return fmt.Sprintf("%s %s rate=%s", l.Carrier, l.Status, rate), nil
	})
	check("Tracking.Get", func() (string, error) {
		if trackingNumber == "" {
			return "SKIP no tracking number", nil
		}
		t, err := c.Tracking.Get(ctx, trackingNumber)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s (%d scans)", t.Carrier, t.Status, len(t.Events)), nil
	})

	fmt.Println("\nQUOTE  (free — no purchase)")
	var rateID string
	check("Rates.List", func() (string, error) {
		r, err := c.Rates.List(ctx, &atoship.RateRequest{FromAddress: from, ToAddress: to, Parcel: parcel})
		if err != nil {
			return "", err
		}
		if len(r) == 0 {
			return "", fmt.Errorf("zero rates returned")
		}
		rateID = r[0].ID
		return fmt.Sprintf("%d rates, cheapest %s %s $%.2f", len(r), r[0].Carrier, r[0].Service, r[0].Rate), nil
	})
	check("Rates.ListBatch", func() (string, error) {
		res, err := c.Rates.ListBatch(ctx, []atoship.BatchRateRequest{
			{ID: "a", RateRequest: atoship.RateRequest{FromAddress: from, ToAddress: to, Parcel: parcel}},
			{ID: "b", RateRequest: atoship.RateRequest{FromAddress: from, ToAddress: to, Parcel: parcel}},
		})
		if err != nil {
			return "", err
		}
		ok := 0
		for _, r := range res {
			if len(r.Rates) > 0 {
				ok++
			}
		}
		return fmt.Sprintf("%d results, %d rated", len(res), ok), nil
	})

	fmt.Println("\nWRITE  (sandbox key — buys nothing)")
	var createdLabel string
	check("Labels.Create", func() (string, error) {
		if rateID == "" {
			return "SKIP no rate id", nil
		}
		l, err := c.Labels.Create(ctx, &atoship.LabelRequest{
			RateID: rateID, FromAddress: from, ToAddress: to, Parcel: parcel,
			Reference: "sdk-verify-" + fmt.Sprint(time.Now().Unix()),
		})
		// A sandbox key validates the request and refuses to execute it. That is
		// a pass: it proves the payload was accepted. Before SandboxError existed
		// this decoded into a zero-valued Label with a nil error, and the harness
		// printed PASS for an object that was never created.
		if atoship.IsSandbox(err) {
			var e *atoship.SandboxError
			errors.As(err, &e)
			return fmt.Sprintf("sandbox validated %d fields", len(e.ValidatedFields)), nil
		}
		if err != nil {
			return "", err
		}
		createdLabel = l.ID
		if l.ID == "" {
			return "", fmt.Errorf("created a label with no id — response did not decode")
		}
		return fmt.Sprintf("%s status=%s", l.ID, l.Status), nil
	})
	check("Labels.Purchase", func() (string, error) {
		if createdLabel == "" {
			return "SKIP sandbox key never created a label to buy", nil
		}
		l, err := c.Labels.Purchase(ctx, createdLabel, nil)
		if atoship.IsSandbox(err) {
			return "sandbox refused to spend", nil
		}
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("status=%s tracking=%q", l.Status, l.TrackingNumber), nil
	})
	check("Labels.BuyBatch", func() (string, error) {
		r, err := c.Labels.BuyBatch(ctx, &atoship.BuyBatchRequest{
			Select: "cheapest",
			Shipments: []atoship.BatchShipment{
				{Reference: "sdkv-1-" + fmt.Sprint(time.Now().Unix()), FromAddress: from, ToAddress: to, Parcel: parcel},
			},
		})
		if atoship.IsSandbox(err) {
			return "sandbox refused to spend", nil
		}
		if err != nil {
			return "", err
		}
		if r.Summary.Total == 0 {
			return "", fmt.Errorf("batch returned an empty manifest — response did not decode")
		}
		return fmt.Sprintf("mode=%s total=%d purchased=%d", r.Mode, r.Summary.Total, r.Summary.Purchased), nil
	})

	fmt.Println("\nERROR HANDLING")
	check("404 -> IsNotFound", func() (string, error) {
		_, err := c.Labels.Get(ctx, "lbl_does_not_exist_zzzz")
		if err == nil {
			return "", fmt.Errorf("expected an error for a missing label, got none")
		}
		if !atoship.IsNotFound(err) {
			return "", fmt.Errorf("expected NOT_FOUND, got %v", err)
		}
		return "typed correctly", nil
	})
	check("bad key -> IsUnauthorized", func() (string, error) {
		bad := atoship.NewClient("ak_test_definitely_invalid", opts...)
		_, err := bad.Account.Get(ctx)
		if err == nil {
			return "", fmt.Errorf("a bogus key was accepted")
		}
		if !atoship.IsUnauthorized(err) {
			return "", fmt.Errorf("expected UNAUTHORIZED, got %v", err)
		}
		return "typed correctly", nil
	})

	fmt.Printf("\n%d passed, %d failed, %d skipped\n", pass, fail, skip)
	if fail > 0 {
		os.Exit(1)
	}
}
