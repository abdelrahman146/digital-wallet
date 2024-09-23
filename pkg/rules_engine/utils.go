package rule_engine

import (
	"strings"
)

// getField extracts the value from a nested field in the JSON-like map using dot notation.
func getField(jsonObject map[string]interface{}, field string) (interface{}, bool) {
	parts := strings.Split(field, ".")
	var current interface{} = jsonObject
	for _, part := range parts {
		if objMap, ok := current.(map[string]interface{}); ok {
			current = objMap[part]
		} else {
			return nil, false
		}
	}
	return current, true
}

func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	default:
		return 0
	}
}
