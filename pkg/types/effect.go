package types

type Effect struct {
	// Type values: "transaction", "tier", "custom"
	Type    string   `json:"type"`
	Formula string   `json:"formula"`
	Params  []string `json:"params"`
}

// effect examples
// reward user with new points based on a calculation made on a given value
// change tiers
