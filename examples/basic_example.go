package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/atoship/go-sdk/atoship"
)

func main() {
	// Initialize the SDK
	client := atoship.NewClient("your-api-key",
		atoship.WithBaseURL("https://api.atoship.com"),
		atoship.WithDebug(true),
	)
	
	ctx := context.Background()
	
	// Example 1: Create an order
	fmt.Println("=== Creating an order ===")
	order, err := client.Orders.Create(ctx, &atoship.CreateOrderRequest{
		OrderNumber:      "GO-ORDER-001",
		RecipientName:    "John Doe",
		RecipientStreet1: "123 Main St",
		RecipientCity:    "San Francisco",
		RecipientState:   "CA",
		RecipientPostal:  "94105",
		RecipientCountry: "US",
		RecipientPhone:   "415-555-0123",
		Items: []atoship.OrderItem{
			{
				Name:       "Go Programming Book",
				SKU:        "BOOK-GO-001",
				Quantity:   2,
				UnitPrice:  29.99,
				Weight:     1.5,
				WeightUnit: "lb",
			},
		},
	})
	
	if err != nil {
		log.Printf("Failed to create order: %v", err)
	} else {
		fmt.Printf("Order created successfully: %s\n", order.ID)
	}
	
	// Example 2: Get shipping rates
	fmt.Println("\n=== Getting shipping rates ===")
	rates, err := client.Shipping.GetRates(ctx, &atoship.RateRequest{
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
	
	if err != nil {
		log.Printf("Failed to get rates: %v", err)
	} else {
		fmt.Printf("Found %d shipping rates:\n", len(rates))
		for _, rate := range rates {
			fmt.Printf("  - %s %s: $%.2f (delivery in %d days)\n",
				rate.Carrier, rate.Service, rate.Rate, rate.DeliveryDays)
		}
	}
	
	// Example 3: Validate an address
	fmt.Println("\n=== Validating an address ===")
	validation, err := client.Addresses.Validate(ctx, &atoship.ValidateAddressRequest{
		Name:       "Jane Smith",
		Street1:    "1600 Amphitheatre Parkway",
		City:       "Mountain View",
		State:      "CA",
		PostalCode: "94043",
		Country:    "US",
	})
	
	if err != nil {
		log.Printf("Failed to validate address: %v", err)
	} else {
		if validation.IsValid {
			fmt.Println("✅ Address is valid")
		} else {
			fmt.Println("❌ Address validation failed")
			if len(validation.Errors) > 0 {
				fmt.Println("Errors:")
				for _, err := range validation.Errors {
					fmt.Printf("  - %s\n", err)
				}
			}
			if len(validation.Suggestions) > 0 {
				fmt.Println("Suggested addresses:")
				for _, addr := range validation.Suggestions {
					fmt.Printf("  - %s, %s, %s %s\n",
						addr.Street1, addr.City, addr.State, addr.PostalCode)
				}
			}
		}
	}
	
	// Example 4: Track a package
	fmt.Println("\n=== Tracking a package ===")
	tracking, err := client.Tracking.Track(ctx, "1Z999AA10123456784")
	
	if err != nil {
		// This will likely fail with a test tracking number
		fmt.Printf("Could not track package: %v\n", err)
	} else {
		fmt.Printf("Package status: %s\n", tracking.Status)
		fmt.Printf("Current location: %s\n", tracking.CurrentLocation)
		if tracking.Delivered {
			fmt.Printf("✅ Package delivered at: %v\n", tracking.ActualDelivery)
		} else if tracking.EstimatedDelivery != nil {
			fmt.Printf("📦 Estimated delivery: %v\n", tracking.EstimatedDelivery)
		}
	}
	
	fmt.Println("\n=== Example completed ===")
}