package validator

import (
	"errors"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"math"
	"regexp"
)

type validatorStruct struct {
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidator() Validator {
	v := &validatorStruct{}
	if err := v.init(); err != nil {
		logger.GetLogger().Panic("failed to initialize validator", logger.Field("error", err))
	}
	return v
}

func (v *validatorStruct) init() error {
	en := en.New()
	v.uni = ut.New(en, en)
	v.trans, _ = v.uni.GetTranslator("en")
	v.validate = validator.New()
	if err := v.validate.RegisterValidation("decimal2", v.Decimal2); err != nil {
		return err
	}
	if err := v.validate.RegisterValidation("slug", v.Slug); err != nil {
		return err
	}
	if err := v.validate.RegisterValidation("alpha", v.Alpha); err != nil {
		return err
	}
	if err := en_translations.RegisterDefaultTranslations(v.validate, v.trans); err != nil {
		return err
	}
	v.registerDecimal2Translation()
	v.RegisterSlugTranslation()
	return v.RegisterAlphaTranslation()
}

func (v *validatorStruct) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

func (v *validatorStruct) GetValidationErrors(err error) map[string]string {
	var errs validator.ValidationErrors
	if ok := errors.As(err, &errs); !ok {
		return nil
	}
	validations := make(map[string]string)
	for _, e := range errs {
		validations[e.Field()] = e.Translate(v.trans)
	}
	return validations
}

func (v *validatorStruct) Decimal2(fl validator.FieldLevel) bool {
	amount, ok := fl.Field().Interface().(float64)
	if !ok {
		return false
	}
	return amount == math.Round(amount*100)/100
}

func (v *validatorStruct) registerDecimal2Translation() error {
	return v.validate.RegisterTranslation("decimal2", v.trans, func(ut ut.Translator) error {
		return ut.Add("decimal2", "{0} must have two digits only", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("decimal2", fe.Field())
		return t
	})
}

func (v *validatorStruct) Slug(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	matched, _ := regexp.MatchString("^[a-z]+(?:[-_][a-z]+)*$", value)
	return matched
}

func (v *validatorStruct) RegisterSlugTranslation() error {
	return v.validate.RegisterTranslation("slug", v.trans, func(ut ut.Translator) error {
		return ut.Add("slug", "{0} must be a valid letters only slug", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("slug", fe.Field())
		return t
	})
}

func (v *validatorStruct) Alpha(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z]+$", value)
	return matched
}

func (v *validatorStruct) RegisterAlphaTranslation() error {
	return v.validate.RegisterTranslation("alpha", v.trans, func(ut ut.Translator) error {
		return ut.Add("alpha", "{0} must have letters only", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alpha", fe.Field())
		return t
	})
}
