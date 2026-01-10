# atoship Go SDK

The official Go SDK for the atoship API. This SDK provides a comprehensive, type-safe interface for all atoship shipping and logistics operations.

## Features

- 🚀 **Type-safe**: Full Go type definitions with compile-time checking
- 🔒 **Secure**: Built-in API key management and request signing
- 🔄 **Robust**: Automatic retries, timeout handling, and error management
- 📦 **Comprehensive**: Covers all atoship API endpoints
- ⚡ **Fast**: Optimized for performance with connection pooling
- 🛡️ **Validated**: Built-in data validation
- 🧪 **Well-tested**: Comprehensive unit tests

## Installation

```bash
go get github.com/atoship/go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/atoship/go-sdk/atoship"
)

func main() {
    // Initialize the SDK
    client := atoship.NewClient("your-api-key")
    
    // Create an order
    order, err := client.Orders.Create(context.Background(), &atoship.CreateOrderRequest{
        OrderNumber:      "ORDER-12345",
        RecipientName:    "John Doe",
        RecipientStreet1: "123 Main St",
        RecipientCity:    "San Francisco",
        RecipientState:   "CA",
        RecipientPostal:  "94105",
        RecipientCountry: "US",
        RecipientPhone:   "415-555-0123",
        Items: []atoship.OrderItem{
            {
                Name:       "Product",
                SKU:        "SKU-001",
                Quantity:   1,
                UnitPrice:  29.99,
                Weight:     2.0,
                WeightUnit: "lb",
            },
        },
    })
    
    if err != nil {
        log.Fatal("Failed to create order:", err)
    }
    
    fmt.Printf("Order created: %s\n", order.ID)
}
```

## SDK Structure

The SDK is organized into logical service groups:

- **Orders**: Order management operations
- **Addresses**: Address validation and management
- **Shipping**: Rate calculation and label generation
- **Tracking**: Package tracking
- **Users**: User and account management
- **Admin**: Administrative operations
- **Carriers**: Carrier-specific operations
- **Webhooks**: Webhook management

## Examples

### Get Shipping Rates

```go
rates, err := client.Shipping.GetRates(context.Background(), &atoship.RateRequest{
    FromAddress: &atoship.Address{
        Street1:    "456 Oak Ave",
        City:       "Los Angeles",
        State:      "CA",
        PostalCode: "90001",
        Country:    "US",
    },
    ToAddress: &atoship.Address{
        Street1:    "789 Pine St",
        City:       "New York",
        State:      "NY",
        PostalCode: "10001",
        Country:    "US",
    },
    Parcel: &atoship.Parcel{
        Length:     10,
        Width:      8,
        Height:     6,
        DimUnit:    "in",
        Weight:     2.5,
        WeightUnit: "lb",
    },
})
```

### Purchase a Label

```go
label, err := client.Shipping.PurchaseLabel(context.Background(), &atoship.PurchaseLabelRequest{
    RateID:  "rate_123456",
    OrderID: "order_789012",
})
```

### Track a Package

```go
tracking, err := client.Tracking.Track(context.Background(), "1Z999AA10123456784")
```

## Error Handling

The SDK provides typed errors for better error handling:

```go
import "github.com/atoship/go-sdk/atoship"

label, err := client.Shipping.PurchaseLabel(ctx, request)
if err != nil {
    if apiErr, ok := err.(*atoship.APIError); ok {
        switch apiErr.Code {
        case atoship.ErrCodeValidation:
            // Handle validation error
        case atoship.ErrCodeRateLimit:
            // Handle rate limit
        case atoship.ErrCodeAuthentication:
            // Handle auth error
        default:
            // Handle other API errors
        }
    }
    // Handle non-API errors
}
```

## Configuration

```go
client := atoship.NewClient("your-api-key",
    atoship.WithBaseURL("https://api.atoship.com"),
    atoship.WithTimeout(30 * time.Second),
    atoship.WithRetryCount(3),
    atoship.WithDebug(true),
)
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -v -cover ./...
```

## Contributing

We welcome contributions! Please see our contributing guidelines for details.

## License

MIT License - see LICENSE file for details.

## Support

- Documentation: https://docs.atoship.com
- API Reference: https://api.atoship.com/docs
- Support: support@atoship.com