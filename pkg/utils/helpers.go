package utils

import (
	"golang.org/x/exp/constraints"
	"strings"
)

type Nullable interface {
	constraints.Ordered | ~*any
}

func Coalesce[T Nullable](options ...T) T {
	var zero T
	for _, option := range options {
		if option != zero {
			return option
		}
	}
	return zero
}

// GetField extracts the value from a nested field in the JSON-like map using dot notation.
func GetField(data map[string]interface{}, field string) (interface{}, bool) {
	parts := strings.Split(field, ".")
	var current interface{} = data
	for _, part := range parts {
		if objMap, ok := current.(map[string]interface{}); ok {
			current = objMap[part]
		} else {
			return nil, false
		}
	}
	return current, true
}
