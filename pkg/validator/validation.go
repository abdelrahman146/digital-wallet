package validator

import (
	"digital-wallet/pkg/logger"
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
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
	return en_translations.RegisterDefaultTranslations(v.validate, v.trans)
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
