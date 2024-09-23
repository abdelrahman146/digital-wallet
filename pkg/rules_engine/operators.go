package rule_engine

import (
	"fmt"
	"github.com/abdelrahman146/digital-wallet/pkg/utils"
	"reflect"
	"regexp"
	"strings"
)

// evaluateOperator handles all comparisons: numeric, string, and date, based on value types.
func evaluateOperator(operator string, fieldValue, ruleValue interface{}) (bool, error) {
	switch fieldVal := fieldValue.(type) {
	case string:
		if utils.IsDate(fieldValue) || utils.IsDate(ruleValue) {
			return evaluateDateOperator(operator, fieldVal, ruleValue)
		}
		return evaluateString(operator, fieldVal, ruleValue.(string))
	case float64:
		return evaluateNumeric(operator, fieldVal, toFloat64(ruleValue))
	case int, int64, int32:
		return evaluateNumeric(operator, float64(reflect.ValueOf(fieldValue).Int()), toFloat64(ruleValue))
	default:
		return false, fmt.Errorf("unsupported value type for comparison: %v", fieldValue)
	}
}

// evaluateDateOperator compares date values using operators like "before", "after", "on".
func evaluateDateOperator(operator string, fieldValue, ruleValue interface{}) (bool, error) {
	fieldDate, err := utils.ParseDate(fieldValue)
	if err != nil {
		return false, fmt.Errorf("field value is not a valid date: %v", err)
	}

	ruleDate, err := utils.ParseDate(ruleValue)
	if err != nil {
		return false, fmt.Errorf("rule value is not a valid date: %v", err)
	}

	switch operator {
	case "==":
		return fieldDate.Equal(ruleDate), nil
	case "!=":
		return !fieldDate.Equal(ruleDate), nil
	case ">":
		return fieldDate.After(ruleDate), nil
	case "<":
		return fieldDate.Before(ruleDate), nil
	case ">=":
		return fieldDate.After(ruleDate) || fieldDate.Equal(ruleDate), nil
	case "<=":
		return fieldDate.Before(ruleDate) || fieldDate.Equal(ruleDate), nil
	default:
		return false, fmt.Errorf("unsupported date operator: %s", operator)
	}
}

// evaluateString handles string comparisons for "==", "!=", ">", "<", ">=", "<="
func evaluateString(operator string, fieldValue string, ruleValue string) (bool, error) {
	switch operator {
	case "==":
		return fieldValue == ruleValue, nil
	case "!=":
		return fieldValue != ruleValue, nil
	case "contains":
		return strings.Contains(fieldValue, ruleValue), nil
	case "in":
		values := strings.Split(ruleValue, ",")
		for _, value := range values {
			if value == fieldValue {
				return true, nil
			}
		}
		return false, nil
	case "notin":
		values := strings.Split(ruleValue, ",")
		for _, value := range values {
			if value == fieldValue {
				return false, nil
			}
		}
		return true, nil
	case "matches":
		return regexp.MatchString(ruleValue, fieldValue)
	default:
		return false, fmt.Errorf("unsupported string operator: %s", operator)
	}
}

// evaluateNumeric handles numeric comparisons for "==", "!=", ">", "<", ">=", "<="
func evaluateNumeric(operator string, fieldValue, ruleValue float64) (bool, error) {
	switch operator {
	case "==":
		return fieldValue == ruleValue, nil
	case "!=":
		return fieldValue != ruleValue, nil
	case ">":
		return fieldValue > ruleValue, nil
	case "<":
		return fieldValue < ruleValue, nil
	case ">=":
		return fieldValue >= ruleValue, nil
	case "<=":
		return fieldValue <= ruleValue, nil
	default:
		return false, fmt.Errorf("unsupported numeric operator: %s", operator)
	}
}
