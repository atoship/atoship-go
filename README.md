# atoship Go SDK

Go client for the [atoship](https://atoship.com) shipping API. Quote a shipment
across carriers, buy the label, track it.

No dependencies outside the standard library.

```bash
go get github.com/atoship/atoship-go
```

```go
import "github.com/atoship/atoship-go/atoship"
```

## Quick start

```go
client := atoship.NewClient(os.Getenv("ATOSHIP_API_KEY"))
ctx := context.Background()

from := atoship.InlineAddress(atoship.Address{
    Name: "Warehouse", Street1: "417 Montgomery St",
    City: "San Francisco", State: "CA", Zip: "94104",
})
to := atoship.InlineAddress(atoship.Address{
    Name: "Jane Doe", Street1: "1600 Pennsylvania Ave NW",
    City: "Washington", State: "DC", Zip: "20500",
})
parcel := atoship.InlineParcel(atoship.Parcel{
    Length: 10, Width: 8, Height: 4, DimensionUnit: "in",
    Weight: 16, WeightUnit: "oz",
})

rates, err := client.Rates.List(ctx, &atoship.RateRequest{
    FromAddress: from, ToAddress: to, Parcel: parcel,
})

label, err := client.Labels.Create(ctx, &atoship.LabelRequest{
    RateID: rates[0].ID, FromAddress: from, ToAddress: to, Parcel: parcel,
})
label, err = client.Labels.Purchase(ctx, label.ID, nil)

fmt.Println(label.TrackingNumber, label.LabelURL)
```

`examples/basic_example.go` is the same flow end to end, including tracking.

## Four things worth knowing before you build on this

**Give the parcel its dimensions.** Carriers bill the greater of actual and
dimensional weight. Quoting a light, bulky box on weight alone returns a price
the carrier will not honour, and the difference arrives later as an adjustment.

**The rate limit is 5 requests per minute, per organization.** Bought one at a
time, a shipment costs three requests — quote, create, purchase — so a loop over
a day's orders stops after the second parcel. That is what `Rates.ListBatch` and
`Labels.BuyBatch` are for: one request, one rate-limit charge, up to 50 quotes or
25 labels.

```go
result, err := client.Labels.BuyBatch(ctx, &atoship.BuyBatchRequest{
    Select: "cheapest",
    Shipments: []atoship.BatchShipment{
        {Reference: "SO-10471", ToAddress: to, Parcel: parcel},
        {Reference: "SO-10472", ToAddress: to2, Parcel: parcel2},
    },
})
```

A nil error does not mean every shipment succeeded. Read `result.Summary` and
each row's `Status`.

**If a batch call times out, do not resubmit.** Read the manifest instead — the
stored rows are authoritative, and a blind resubmit can buy a second label:

```go
result, err := client.Labels.GetBatch(ctx, batchID)
```

**`Status: "needs_review"` on a row means the wallet was charged but the label
was not recorded.** That shipment is held for reconciliation. Do not re-buy it.

## Errors

Every failure is an `*atoship.Error` carrying the API's own code. Branch on the
code, not the message:

```go
label, err := client.Labels.Purchase(ctx, id, nil)

var apiErr *atoship.Error
if errors.As(err, &apiErr) && apiErr.Code == atoship.ErrInsufficientBalance {
    // top up, or enable auto-recharge
}
```

`IsNotFound`, `IsRateLimited`, `IsUnauthorized` and `IsSandbox` cover the common
checks.

## Sandbox

A key beginning `ak_test_` validates requests and refuses to spend money.
Anything that would buy returns an `*atoship.SandboxError` — the request shape
was accepted, only the key needs swapping:

```go
_, err := client.Labels.Create(ctx, req)
if atoship.IsSandbox(err) {
    // request is well-formed; use an ak_live_ key to actually create it
}
```

Reads and quotes work normally on a sandbox key.

## What is covered

| Service | Calls |
|---|---|
| `Account` | Get |
| `Rates` | List, ListBatch |
| `Labels` | Create, Purchase, Get, List, Void, DownloadURL, BuyBatch, GetBatch |
| `Addresses` | Create, List, Validate, Verify |
| `Tracking` | Get |
| `CarrierAccounts` | List, Get, Create, Update, Delete |
| `Webhooks` | List, Get, Create, Update, Patch, Delete |
| `Returns` | Create, List |
| `Insurances` | Create, Get, List, FileClaim, ListClaims, CancelClaim |
| `Pickups` | Create, Get, List, Cancel |
| `Orders` | GetByPlatform |

Client options: `WithBaseURL`, `WithTimeout`, `WithHTTPClient`, `WithUserAgent`.

## Verifying against the live API

`verify/` calls every method above against the real API and prints a pass/fail
line per endpoint. Nothing in it is stubbed, which is the point:

```bash
ATOSHIP_API_KEY=ak_test_… go run ./verify
```

Use a sandbox key. Reads and quotes run for real; writes come back as
`SandboxError` without buying anything.

## Documentation

- API reference: https://atoship.com/docs/api-reference
- Issues: https://github.com/atoship/atoship-go/issues

## License

MIT
