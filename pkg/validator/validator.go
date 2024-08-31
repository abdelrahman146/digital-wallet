package validator

type Validator interface {
	ValidateStruct(s interface{}) error
	GetValidationErrors(err error) map[string]string
}

var v Validator

func GetValidator() Validator {
	if v == nil {
		v = NewValidator()
	}
	return v
}
