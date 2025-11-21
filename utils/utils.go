package utils

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
)

func ParseExpirationDate(expirationDateString string) (time.Time, error) {
	dateTime, err := dateparse.ParseAny(expirationDateString)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %w", err)
	}
	return dateTime, nil
}
