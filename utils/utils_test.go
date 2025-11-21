package utils

import (
	"testing"
)

func TestParseExpirationDate(t *testing.T) {
	expirations := []string{"2021-01-01", "2021-01-01 12:00:00"}

	for _, expiration := range expirations {
		t.Run(expiration, func(t *testing.T) {
			_, err := ParseExpirationDate(expiration)
			if err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}
		})
	}

	badExpirations := []string{"2021-0111-01 12:00:00", "20211-01 12:00:00.000000"}
	for _, expiration := range badExpirations {
		t.Run(expiration, func(t *testing.T) {
			_, err := ParseExpirationDate(expiration)
			if err == nil {
				t.Errorf("Expected an error, but got none")
			}
		})
	}
}
