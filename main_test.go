package main

import (
	"context"
	"fmt"
	"testing"
)

func TestDomainAvailability(t *testing.T) {
	domains := []map[string]interface{}{
		{
			"domain":              "google.com",
			"should_be_available": false,
		},
		{
			"domain":              "ajkladhsjkfayfasdfbhjkbcvghajkdshcjkhauijdshfjkasdhfjkl.com",
			"should_be_available": true,
		},
	}

	for _, domain := range domains {
		fmt.Printf("Checking: %s\n", domain)

		available := checkDomainAvailabilityWhoIs(domain["domain"].(string))

		if available != domain["should_be_available"].(bool) {
			t.Errorf("Expected domain %s to be %v, but got %v", domain["domain"].(string), domain["should_be_available"].(bool), available)
		}
	}
}

func TestDomainTLD(t *testing.T) {
	domains := []map[string]interface{}{
		{
			"domain":              "google",
			"should_be_available": false,
		},
	}

	for _, domain := range domains {
		fmt.Printf("Checking: %s\n", domain["domain"].(string))

		_, output, err := CheckAvailableDomainTLD(context.Background(), nil, AvailableDomainTLDInput{
			Domain:      domain["domain"].(string),
			OnlyPopular: true,
			OnlyCountry: false,
		})

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if len(output.AvailableDomains) != 0 {
			t.Errorf("Expected available domains to be 0, but got %v", len(output.AvailableDomains))
		}
	}
}
