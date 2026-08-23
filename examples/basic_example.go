// Quote a shipment, buy the cheapest rate, print the label, then track it.
//
//	export ATOSHIP_API_KEY=ak_test_…
//	go run ./examples
//
// A key beginning ak_test_ is a sandbox key: everything below runs, nothing is
// bought and nothing is charged.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/atoship/atoship-go/atoship"
)

func main() {
	key := os.Getenv("ATOSHIP_API_KEY")
	if key == "" {
		log.Fatal("set ATOSHIP_API_KEY")
	}

	client := atoship.NewClient(key, atoship.WithUserAgent("atoship-example/1"))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Cheapest way to check the key works: no arguments, buys nothing.
	acct, err := client.Account.Get(ctx)
	if err != nil {
		if atoship.IsUnauthorized(err) {
			log.Fatal("that API key was rejected")
		}
		log.Fatalf("account: %v", err)
	}
	fmt.Printf("account %s — $%.2f available\n", acct.Email, acct.Balance.Available)

	from := atoship.InlineAddress(atoship.Address{
		Name: "Warehouse", Street1: "417 Montgomery St", City: "San Francisco",
		State: "CA", Zip: "94104", Country: "US",
	})
	to := atoship.InlineAddress(atoship.Address{
		Name: "Jane Doe", Street1: "1600 Pennsylvania Ave NW", City: "Washington",
		State: "DC", Zip: "20500", Country: "US",
	})
	// Dimensions are not optional in practice: carriers bill the greater of
	// actual and dimensional weight, so leaving them off a light bulky box
	// quotes a price the carrier will not honour.
	parcel := atoship.InlineParcel(atoship.Parcel{
		Length: 10, Width: 8, Height: 4, DimensionUnit: "in",
		Weight: 16, WeightUnit: "oz",
	})

	rates, err := client.Rates.List(ctx, &atoship.RateRequest{
		FromAddress: from, ToAddress: to, Parcel: parcel,
	})
	if err != nil {
		log.Fatalf("rates: %v", err)
	}
	if len(rates) == 0 {
		log.Fatal("no rates for that shipment")
	}
	fmt.Printf("\n%d rates:\n", len(rates))
	for _, r := range rates {
		fmt.Printf("  %-10s %-24s $%6.2f  %d days\n", r.Carrier, r.Service, r.Rate, r.DeliveryDays)
	}

	// Rates come back cheapest first. Buying by RateID buys the exact price you
	// just showed someone; naming a carrier and service re-prices at purchase.
	cheapest := rates[0]
	fmt.Printf("\nbuying %s %s at $%.2f\n", cheapest.Carrier, cheapest.Service, cheapest.Rate)

	label, err := client.Labels.Create(ctx, &atoship.LabelRequest{
		RateID:      cheapest.ID,
		FromAddress: from, ToAddress: to, Parcel: parcel,
		Reference: "example-001",
	})
	if err != nil {
		log.Fatalf("create label: %v", err)
	}

	label, err = client.Labels.Purchase(ctx, label.ID, nil)
	if err != nil {
		var apiErr *atoship.Error
		if errors.As(err, &apiErr) && apiErr.Code == atoship.ErrInsufficientBalance {
			log.Fatal("not enough balance — top up or enable auto-recharge")
		}
		log.Fatalf("purchase: %v", err)
	}
	fmt.Printf("label %s\n  pdf: %s\n", label.ID, label.LabelURL)

	// On some services the carrier assigns the number after the buy. An empty
	// tracking number with TrackingPending set is not a failure.
	if label.TrackingNumber == "" && label.TrackingPending {
		fmt.Println("  tracking number still being assigned by the carrier")
		return
	}
	fmt.Printf("  tracking: %s\n", label.TrackingNumber)

	tracking, err := client.Tracking.Get(ctx, label.TrackingNumber)
	if err != nil {
		log.Fatalf("tracking: %v", err)
	}
	// A label bought seconds ago has no scans yet. pre_transit with no events is
	// the expected state, not an error.
	fmt.Printf("  status: %s (%d scans)\n", tracking.Status, len(tracking.Events))
}
