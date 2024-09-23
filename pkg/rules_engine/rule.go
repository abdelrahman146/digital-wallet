package rule_engine

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Rule represents a validation rule. It can include nested rules for AND/OR/NOT logic.
type Rule struct {
	Field    string      `json:"field,omitempty"`    // Field name to apply the rule to (can use dot notation for nested fields like "address.city")
	Operator string      `json:"operator,omitempty"` // Operator for the rule (e.g., "==", ">", "before", "after", etc.)
	Val      interface{} `json:"value,omitempty"`    // Val to compare the field with
	Logic    string      `json:"logic,omitempty"`    // Logical operator for combining rules ("AND", "OR", "NOT")
	Rules    []Rule      `json:"rules,omitempty"`    // Nested rules (used with logic operators)
}

// Value Marshal
func (a Rule) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Rule) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
