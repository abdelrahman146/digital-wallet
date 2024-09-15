package utils

import (
	"golang.org/x/exp/constraints"
	"time"
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

func GetDateAfter(interval time.Duration) time.Time {
	return time.Now().Add(interval)
}
