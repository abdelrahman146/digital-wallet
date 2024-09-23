package utils

import (
	"errors"
	"time"
)

func GetDateAfter(interval time.Duration) time.Time {
	return time.Now().Add(interval)
}

// ParseDate converts a date string into a time.Time object. Supports multiple date formats.
func ParseDate(value interface{}) (time.Time, error) {
	strValue, ok := value.(string)
	if !ok {
		return time.Time{}, errors.New("invalid date format")
	}

	layouts := []string{
		time.RFC3339,       // "2006-01-02T15:04:05Z07:00"
		"2006-01-02",       // "2006-01-02"
		"2006-01-02 15:04", // "2006-01-02 15:04"
		"02-01-2006",       // "02-01-2006" (European format)
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, strValue)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, errors.New("could not parse date")
}

// IsDate checks if a value is a date string.
func IsDate(value interface{}) bool {
	_, err := ParseDate(value)
	return err == nil
}
