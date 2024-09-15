package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB Interface for JSONB Field
type JSONB map[string]interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (a *JSONB) StructToJSONB(value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		return err
	}
	return nil
}
