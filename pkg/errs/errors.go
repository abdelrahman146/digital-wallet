package errs

import "digital-wallet/pkg/utils"

func NewValidationError(text string, code string, fields map[string]string) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Could not understand the request due to invalid syntax."),
		Code:     utils.Coalesce(code, "VALIDATION_ERROR"),
		HttpCode: 400,
		Fields:   fields,
	}
}

func NewNotFoundError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Resource could not be found"),
		Code:     utils.Coalesce(code, "NOT_FOUND"),
		HttpCode: 404,
		Original: err,
	}
}

func NewConflictError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Uncompleted request due to a conflict with the current state of the target resources."),
		Code:     utils.Coalesce(code, "CONFLICT"),
		HttpCode: 409,
		Original: err,
	}
}

func NewUnauthorizedError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Authorization required for the target resources."),
		Code:     utils.Coalesce(code, "UNAUTHORIZED"),
		HttpCode: 401,
		Original: err,
	}
}

func NewForbiddenError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Request understood, but refused"),
		Code:     utils.Coalesce(code, "FORBIDDEN"),
		HttpCode: 403,
		Original: err,
	}
}

func NewBadRequestError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Could not understand the request due to invalid syntax."),
		Code:     utils.Coalesce(code, "BAD_REQUEST"),
		HttpCode: 400,
		Original: err,
	}
}

func NewPaymentRequiredError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Payment required for the target resources."),
		Code:     utils.Coalesce(code, "PAYMENT_REQUIRED"),
		HttpCode: 402,
		Original: err,
	}
}
func NewMethodNotAllowedError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Request method not supported for the target resources."),
		Code:     utils.Coalesce(code, "METHOD_NOT_ALLOWED"),
		HttpCode: 405,
		Original: err,
	}
}

func NewNotAcceptableError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Request not acceptable for the target resources."),
		Code:     utils.Coalesce(code, "NOT_ACCEPTABLE"),
		HttpCode: 406,
		Original: err,
	}
}
func NewTimeoutError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Request timeout"),
		Code:     utils.Coalesce(code, "TIMEOUT"),
		HttpCode: 408,
		Original: err,
	}
}

func NewPayloadTooLargeError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Request payload too large"),
		Code:     utils.Coalesce(code, "PAYLOAD_TOO_LARGE"),
		HttpCode: 413,
		Original: err,
	}
}

func NewUnsupportedMediaTypeError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Unsupported media type"),
		Code:     utils.Coalesce(code, "UNSUPPORTED_MEDIA_TYPE"),
		HttpCode: 415,
		Original: err,
	}
}

func NewExpectationFailedError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Expectation failed"),
		Code:     utils.Coalesce(code, "EXPECTATION_FAILED"),
		HttpCode: 417,
		Original: err,
	}
}

func NewTooManyRequestsError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Too many requests"),
		Code:     utils.Coalesce(code, "TOO_MANY_REQUESTS"),
		HttpCode: 429,
		Original: err,
	}
}

func NewUnprocessableEntityError(text string, code string, err error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Well formed request, but was unable to be followed due to semantic errors."),
		Code:     utils.Coalesce(code, "UNPROCESSABLE_ENTITY"),
		HttpCode: 422,
		Original: err,
	}
}

func NewInternalError(text string, code string, e error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "Encountered an unexpected condition that prevented it from fulfilling the request."),
		Code:     utils.Coalesce(code, "INTERNAL_ERROR"),
		HttpCode: 500,
		Original: e,
	}
}
func NewNotImplementedError(text string, code string, e error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "The server does not support the functionality required to fulfill the request."),
		Code:     utils.Coalesce(code, "NOT_IMPLEMENTED"),
		HttpCode: 501,
		Original: e,
	}
}

func NewServiceUnavailableError(text string, code string, e error) CustomError {
	return CustomError{
		Message:  utils.Coalesce(text, "The server is currently unavailable."),
		Code:     utils.Coalesce(code, "SERVICE_UNAVAILABLE"),
		HttpCode: 503,
		Original: e,
	}
}

func NewUnknownError(e error) CustomError {
	return CustomError{
		Message:  "Unhandled Error",
		Code:     "UNKNOWN_ERROR",
		HttpCode: 520,
		Original: e,
	}
}
