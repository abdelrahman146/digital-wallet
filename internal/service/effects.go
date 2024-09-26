package service

import (
	"context"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/utils"
)

func ApplyEffect(ctx context.Context, program model.Program, data map[string]interface{}) error {
	switch program.Effect["type"] {
	case "FIXED":
		return EvaluateFixedEffect(ctx, program, data)
	case "FORMULA":
		return EvaluateFormulaEffect(ctx, program, data)
	case "PROMOTE":
		return EvaluateTierEffect(ctx, program, data)
	case "CALL":
		return EvaluateCallEffect(ctx, program, data)
	default:
		return errs.NewUnprocessableEntityError("Program has Invalid Effect Type", "PROGRAM_INVALID_EFFECT_TYPE", nil)
	}
}

func EvaluateFixedEffect(ctx context.Context, program model.Program, data map[string]interface{}) error {
	return nil
}

func EvaluateFormulaEffect(ctx context.Context, program model.Program, data map[string]interface{}) error {
	formula, ok := program.Effect["formula"].(string)
	if !ok || formula == "" {
		return errs.NewUnprocessableEntityError("Program type is 'FORMULA' but doesn't have a formula", "PROGRAM_INVALID_EFFECT_FORMULA", nil)
	}
	params, ok := program.Effect["parameters"].([]string)
	if !ok || params == nil {
		return errs.NewUnprocessableEntityError("Program type is 'FORMULA' but doesn't have parameters", "PROGRAM_INVALID_EFFECT_PARAMETERS", nil)
	}
	paramValues := make(map[string]interface{}, len(params))
	for _, param := range params {
		paramValues[param], ok = utils.GetField(data, param)
		if !ok {
			return errs.NewUnprocessableEntityError(fmt.Sprintf("Invalid Pramater: %s", param), "PROGRAM_INVALID_EFFECT_PRAMATER_VALUE", nil)
		}
	}
	exp, err := govaluate.NewEvaluableExpression(formula)
	if err != nil {
		return errs.NewUnprocessableEntityError("Invalid Formula: "+formula, "PROGRAM_INVALID_EFFECT_FORMULA", err)
	}
	_, err = exp.Evaluate(paramValues)
	if err != nil {
		return errs.NewUnprocessableEntityError(fmt.Sprintf("Unable to evaluate formula: %s and params %v", formula, params), "PROGRAM_INVALID_EFFECT_FORMULA", err)
	}
	// TODO: continue
	return nil
}

func EvaluateCallEffect(ctx context.Context, program model.Program, data map[string]interface{}) error {
	return nil
}

func EvaluateTierEffect(ctx context.Context, program model.Program, data map[string]interface{}) error {
	return nil
}
