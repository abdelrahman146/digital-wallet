package rule_engine

import (
	"encoding/json"
	"fmt"
)

func Evaluate(rule Rule, data []byte) (bool, error) {
	var jsonObject map[string]interface{}
	err := json.Unmarshal(data, &jsonObject)
	if err != nil {
		return false, fmt.Errorf("invalid JSON data: %v", err)
	}
	return evaluateRule(rule, jsonObject)
}

// evaluateRule recursively evaluates rules, including logical combinations (AND, OR, NOT).
func evaluateRule(rule Rule, jsonObject map[string]interface{}) (bool, error) {
	if rule.Logic != "" {
		return evaluateLogic(rule, jsonObject)
	}

	fieldValue, exists := getField(jsonObject, rule.Field)
	if !exists {
		return false, fmt.Errorf("field %s not found", rule.Field)
	}

	// Handle array fields - Check if it's an array and operator is related to array handling (e.g., "any", "all")
	if array, ok := fieldValue.([]interface{}); ok && (rule.Operator == "any" || rule.Operator == "all") {
		return evaluateArray(rule, array, rule.Operator == "all")
	}

	// Evaluate the operator
	return evaluateOperator(rule.Operator, fieldValue, rule.Val)
}

// evaluateArray applies a rule to each element in an array (supports "all" or "any").
func evaluateArray(rule Rule, jsonArray []interface{}, matchAll bool) (bool, error) {
	for _, element := range jsonArray {
		jsonObject, ok := element.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("array element is not a JSON object")
		}

		// Iterate through all sub-rules and evaluate them for each array element
		for _, subRule := range rule.Rules {
			result, err := evaluateRule(subRule, jsonObject)
			if err != nil {
				return false, err
			}

			// For "all" case, return false if any rule fails for any element
			if matchAll && !result {
				return false, nil
			}

			// For "any" case, return true if any rule passes for at least one element
			if !matchAll && result {
				return true, nil
			}
		}
	}

	// If it's an "all" rule, return true because all elements passed all rules
	// If it's an "any" rule, return false because no element satisfied any rule
	return matchAll, nil
}

// evaluateLogic handles AND, OR, and NOT operations between rules.
func evaluateLogic(rule Rule, jsonObject map[string]interface{}) (bool, error) {
	switch rule.Logic {
	case "AND":
		for _, subRule := range rule.Rules {
			result, err := evaluateRule(subRule, jsonObject)
			if err != nil || !result {
				return false, err
			}
		}
		return true, nil
	case "OR":
		for _, subRule := range rule.Rules {
			result, err := evaluateRule(subRule, jsonObject)
			if result && err == nil {
				return true, nil
			}
		}
		return false, nil
	case "NOT":
		if len(rule.Rules) != 1 {
			return false, fmt.Errorf("NOT logic must have exactly one sub-rule")
		}
		result, err := evaluateRule(rule.Rules[0], jsonObject)
		return !result, err
	default:
		return false, fmt.Errorf("unsupported logic operator: %s", rule.Logic)
	}
}
