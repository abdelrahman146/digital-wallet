package rule_engine

import (
	"testing"
)

// Utility function to make it easier to handle test cases
func testEvaluate(t *testing.T, rule Rule, data string, expected bool) {
	result, err := Evaluate(rule, []byte(data))
	if err != nil {
		t.Fatalf("Evaluation failed with error: %v", err)
	}
	if result != expected {
		t.Fatalf("Expected %v but got %v for rule: %+v and data: %s", expected, result, rule, data)
	}
}

func TestEvaluateSimpleStringEquality(t *testing.T) {
	rule := Rule{
		Field:    "name",
		Operator: "==",
		Val:      "John",
	}

	data := `{"name": "John"}`
	testEvaluate(t, rule, data, true)

	data = `{"name": "Doe"}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateStringContains(t *testing.T) {
	rule := Rule{
		Field:    "description",
		Operator: "contains",
		Val:      "world",
	}

	data := `{"description": "Hello world"}`
	testEvaluate(t, rule, data, true)

	data = `{"description": "Hello everyone"}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateStringIn(t *testing.T) {
	rule := Rule{
		Field:    "status",
		Operator: "in",
		Val:      "pending,approved,rejected",
	}

	data := `{"status": "approved"}`
	testEvaluate(t, rule, data, true)

	data = `{"status": "archived"}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateStringNotIn(t *testing.T) {
	rule := Rule{
		Field:    "status",
		Operator: "notin",
		Val:      "pending,approved,rejected",
	}

	data := `{"status": "archived"}`
	testEvaluate(t, rule, data, true)

	data = `{"status": "approved"}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateStringMatches(t *testing.T) {
	rule := Rule{
		Field:    "email",
		Operator: "matches",
		Val:      `^\S+@\S+\.\S+$`, // simple regex for email
	}

	data := `{"email": "test@example.com"}`
	testEvaluate(t, rule, data, true)

	data = `{"email": "invalid-email"}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateNumericGreaterThan(t *testing.T) {
	rule := Rule{
		Field:    "age",
		Operator: ">",
		Val:      18,
	}

	data := `{"age": 20}`
	testEvaluate(t, rule, data, true)

	data = `{"age": 15}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateDateBefore(t *testing.T) {
	rule := Rule{
		Field:    "birthdate",
		Operator: "<",
		Val:      "2000-01-01",
	}

	data := `{"birthdate": "1995-05-12"}`
	testEvaluate(t, rule, data, true)

	data = `{"birthdate": "2005-01-01"}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateAndLogic(t *testing.T) {
	rule := Rule{
		Logic: "AND",
		Rules: []Rule{
			{
				Field:    "age",
				Operator: ">",
				Val:      18,
			},
			{
				Field:    "name",
				Operator: "==",
				Val:      "John",
			},
		},
	}

	data := `{"name": "John", "age": 25}`
	testEvaluate(t, rule, data, true)

	data = `{"name": "Doe", "age": 25}`
	testEvaluate(t, rule, data, false)

	data = `{"name": "John", "age": 15}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateOrLogic(t *testing.T) {
	rule := Rule{
		Logic: "OR",
		Rules: []Rule{
			{
				Field:    "age",
				Operator: ">",
				Val:      18,
			},
			{
				Field:    "name",
				Operator: "==",
				Val:      "John",
			},
		},
	}

	data := `{"name": "John", "age": 15}`
	testEvaluate(t, rule, data, true)

	data = `{"name": "Doe", "age": 25}`
	testEvaluate(t, rule, data, true)

	data = `{"name": "Doe", "age": 15}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateNotLogic(t *testing.T) {
	rule := Rule{
		Logic: "NOT",
		Rules: []Rule{
			{
				Field:    "age",
				Operator: ">",
				Val:      18,
			},
		},
	}

	data := `{"age": 20}`
	testEvaluate(t, rule, data, false)

	data = `{"age": 15}`
	testEvaluate(t, rule, data, true)
}

func TestEvaluateNestedFields(t *testing.T) {
	rule := Rule{
		Field:    "address.city",
		Operator: "==",
		Val:      "New York",
	}

	data := `{"address": {"city": "New York"}}`
	testEvaluate(t, rule, data, true)

	data = `{"address": {"city": "Los Angeles"}}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateArrayAny(t *testing.T) {
	rule := Rule{
		Field:    "family",
		Operator: "any",
		Rules: []Rule{
			{
				Field:    "age",
				Operator: ">",
				Val:      18,
			},
		},
	}

	data := `{"family": [{"name": "John", "age": 15}, {"name": "Jane", "age": 25}]}`
	testEvaluate(t, rule, data, true)

	data = `{"family": [{"name": "John", "age": 15}, {"name": "Jane", "age": 17}]}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateArrayAll(t *testing.T) {
	rule := Rule{
		Field:    "family",
		Operator: "all",
		Rules: []Rule{
			{
				Field:    "age",
				Operator: ">",
				Val:      18,
			},
		},
	}

	data := `{"family": [{"name": "John", "age": 25}, {"name": "Jane", "age": 30}]}`
	testEvaluate(t, rule, data, true)

	data = `{"family": [{"name": "John", "age": 25}, {"name": "Jane", "age": 17}]}`
	testEvaluate(t, rule, data, false)
}

func TestEvaluateArrayMultipleRulesAny(t *testing.T) {
	rule := Rule{
		Field:    "family",
		Operator: "any",
		Rules: []Rule{
			{
				Field:    "age",
				Operator: ">",
				Val:      18,
			},
			{
				Field:    "name",
				Operator: "==",
				Val:      "John",
			},
		},
	}

	data := `{"family": [{"name": "John", "age": 15}, {"name": "Jane", "age": 25}]}`
	testEvaluate(t, rule, data, true) // One rule (name == "John") is satisfied

	data = `{"family": [{"name": "John", "age": 15}, {"name": "Jane", "age": 17}]}`
	testEvaluate(t, rule, data, true) // One rule (name == "John") is satisfied

	data = `{"family": [{"name": "Doe", "age": 15}, {"name": "Jane", "age": 17}]}`
	testEvaluate(t, rule, data, false) // No rule is satisfied
}

func TestEvaluateArrayMultipleRulesAll(t *testing.T) {
	rule := Rule{
		Field:    "family",
		Operator: "all",
		Rules: []Rule{
			{
				Field:    "age",
				Operator: ">",
				Val:      18,
			},
			{
				Field:    "name",
				Operator: "==",
				Val:      "John",
			},
		},
	}

	data := `{"family": [{"name": "John", "age": 25}, {"name": "John", "age": 30}]}`
	testEvaluate(t, rule, data, true) // All elements satisfy all rules

	data = `{"family": [{"name": "John", "age": 25}, {"name": "Jane", "age": 30}]}`
	testEvaluate(t, rule, data, false) // Second element fails the "name == John" rule
}

func TestComplexRuleWithAndOrNot(t *testing.T) {
	rule := Rule{
		Logic: "AND",
		Rules: []Rule{
			{
				Field:    "user.name",
				Operator: "==",
				Val:      "John",
			},
			{
				Logic: "OR",
				Rules: []Rule{
					{
						Field:    "user.age",
						Operator: ">",
						Val:      30,
					},
					{
						Field:    "user.occupation",
						Operator: "in",
						Val:      "engineer,doctor,teacher",
					},
				},
			},
			{
				Logic: "NOT",
				Rules: []Rule{
					{
						Field:    "user.location",
						Operator: "==",
						Val:      "New York",
					},
				},
			},
		},
	}

	// Data that should pass all rules
	data := `{
		"user": {
			"name": "John",
			"age": 35,
			"occupation": "doctor",
			"location": "California"
		}
	}`
	testEvaluate(t, rule, data, true)

	// Data that fails the NOT rule
	data = `{
		"user": {
			"name": "John",
			"age": 35,
			"occupation": "engineer",
			"location": "New York"
		}
	}`
	testEvaluate(t, rule, data, false) // Should fail due to NOT on "location == New York"

	// Data that fails the OR condition
	data = `{
		"user": {
			"name": "John",
			"age": 25,
			"occupation": "student",
			"location": "California"
		}
	}`
	testEvaluate(t, rule, data, false) // Fails OR (age <= 30 and occupation not in "engineer,doctor,teacher")

	// Data that fails the AND condition
	data = `{
		"user": {
			"name": "Doe",
			"age": 40,
			"occupation": "teacher",
			"location": "California"
		}
	}`
	testEvaluate(t, rule, data, false) // Fails because "name" is not "John"
}

func TestArrayHandlingWithMultipleNestedRules(t *testing.T) {
	rule := Rule{
		Field:    "family",
		Operator: "all",
		Rules: []Rule{
			{
				Logic: "AND",
				Rules: []Rule{
					{
						Field:    "age",
						Operator: ">",
						Val:      18,
					},
					{
						Field:    "name",
						Operator: "matches",
						Val:      "^[A-Za-z]+$", // Only alphabet names
					},
					{
						Field:    "relationship",
						Operator: "in",
						Val:      "father,mother,sibling",
					},
				},
			},
		},
	}

	// Data where all family members meet the criteria
	data := `{
		"family": [
			{"name": "John", "age": 45, "relationship": "father"},
			{"name": "Jane", "age": 42, "relationship": "mother"},
			{"name": "Jake", "age": 19, "relationship": "sibling"}
		]
	}`
	testEvaluate(t, rule, data, true)

	// Data where one family member fails the "relationship in" rule
	data = `{
		"family": [
			{"name": "John", "age": 45, "relationship": "father"},
			{"name": "Jane", "age": 42, "relationship": "mother"},
			{"name": "Jake", "age": 19, "relationship": "cousin"}
		]
	}`
	testEvaluate(t, rule, data, false) // Fails because "cousin" is not in "father,mother,sibling"

	// Data where one family member fails the regex match for name
	data = `{
		"family": [
			{"name": "John", "age": 45, "relationship": "father"},
			{"name": "Jane", "age": 42, "relationship": "mother"},
			{"name": "Jake123", "age": 19, "relationship": "sibling"}
		]
	}`
	testEvaluate(t, rule, data, false) // Fails because "Jake123" doesn't match the regex

	// Data where one family member fails the age rule
	data = `{
		"family": [
			{"name": "John", "age": 45, "relationship": "father"},
			{"name": "Jane", "age": 16, "relationship": "sibling"},
			{"name": "Jake", "age": 19, "relationship": "sibling"}
		]
	}`
	testEvaluate(t, rule, data, false) // Fails because Jane is 16 (< 18)
}

func TestComplexDateAndStringOperations(t *testing.T) {
	rule := Rule{
		Logic: "AND",
		Rules: []Rule{
			{
				Field:    "user.registration_date",
				Operator: "<=",
				Val:      "2023-01-01",
			},
			{
				Field:    "user.email",
				Operator: "matches",
				Val:      `^\S+@\S+\.\S+$`, // Simple email regex
			},
			{
				Field:    "user.status",
				Operator: "in",
				Val:      "active,pending",
			},
		},
	}

	// Data that meets all the criteria
	data := `{
		"user": {
			"registration_date": "2022-12-01",
			"email": "john@example.com",
			"status": "active"
		}
	}`
	testEvaluate(t, rule, data, true)

	// Data that fails the "before" date rule
	data = `{
		"user": {
			"registration_date": "2023-05-01",
			"email": "john@example.com",
			"status": "active"
		}
	}`
	testEvaluate(t, rule, data, false) // Fails date rule

	// Data that fails the email regex match
	data = `{
		"user": {
			"registration_date": "2022-12-01",
			"email": "invalid-email",
			"status": "active"
		}
	}`
	testEvaluate(t, rule, data, false) // Fails email regex

	// Data that fails the "in" rule for status
	data = `{
		"user": {
			"registration_date": "2022-12-01",
			"email": "john@example.com",
			"status": "suspended"
		}
	}`
	testEvaluate(t, rule, data, false) // Fails "status in" rule
}
