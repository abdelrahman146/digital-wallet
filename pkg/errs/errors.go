package errs

func NewValidationError(text string, fields map[string]string) CustomError {
	return CustomError{
		Message:  text,
		Code:     "VALIDATION_ERROR",
		Desc:     "Could not understand the request due to invalid syntax.",
		HttpCode: 400,
		Fields:   fields,
	}
}

func NewNotFoundError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "NOT_FOUND",
		Desc:     "Resource could not be found",
		HttpCode: 404,
		Original: err,
	}
}

func NewConflictError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "CONFLICT",
		Desc:     "Uncompleted request due to a conflict with the current state of the target resources.",
		HttpCode: 409,
		Original: err,
	}
}

func NewUnauthorizedError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "UNAUTHORIZED",
		Desc:     "Authentication required for the target resources.",
		HttpCode: 401,
		Original: err,
	}
}

func NewForbiddenError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "FORBIDDEN",
		Desc:     "Request understood, but refused",
		HttpCode: 403,
		Original: err,
	}
}

func NewBadRequestError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "BAD_REQUEST",
		Desc:     "Could not understand the request due to invalid syntax.",
		HttpCode: 400,
		Original: err,
	}
}

func NewPaymentRequiredError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "PAYMENT_REQUIRED",
		Desc:     "Payment required for the target resources.",
		HttpCode: 402,
		Original: err,
	}
}
func NewMethodNotAllowedError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "METHOD_NOT_ALLOWED",
		Desc:     "Request method not supported for the target resources.",
		HttpCode: 405,
		Original: err,
	}
}

func NewNotAcceptableError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "NOT_ACCEPTABLE",
		Desc:     "Request not acceptable for the target resources.",
		HttpCode: 406,
		Original: err,
	}
}
func NewTimeoutError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "TIMEOUT",
		Desc:     "Request timeout",
		HttpCode: 408,
		Original: err,
	}
}

func NewPayloadTooLargeError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "PAYLOAD_TOO_LARGE",
		Desc:     "Request payload too large",
		HttpCode: 413,
		Original: err,
	}
}

func NewUnsupportedMediaTypeError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "UNSUPPORTED_MEDIA_TYPE",
		Desc:     "Unsupported media type",
		HttpCode: 415,
		Original: err,
	}
}

func NewExpectationFailedError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "EXPECTATION_FAILED",
		Desc:     "Expectation failed",
		HttpCode: 417,
		Original: err,
	}
}

func NewTooManyRequestsError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "TOO_MANY_REQUESTS",
		Desc:     "Too many requests",
		HttpCode: 429,
		Original: err,
	}
}

func NewUnprocessableEntityError(text string, err error) CustomError {
	return CustomError{
		Message:  text,
		Code:     "Well formed request, but was unable to be followed due to semantic errors.",
		Desc:     "Unprocessable entity",
		HttpCode: 422,
		Original: err,
	}
}

func NewInternalError(text string, e error) CustomError {
	return CustomError{
		Message:  text,
		Desc:     "Encountered an unexpected condition that prevented it from fulfilling the request.",
		Code:     "INTERNAL_ERROR",
		HttpCode: 500,
		Original: e,
	}
}
func NewNotImplementedError(text string, e error) CustomError {
	return CustomError{
		Message:  text,
		Desc:     "The server does not support the functionality required to fulfill the request.",
		Code:     "NOT_IMPLEMENTED",
		HttpCode: 501,
		Original: e,
	}
}

func NewServiceUnavailableError(text string, e error) CustomError {
	return CustomError{
		Message:  text,
		Desc:     "The server is currently unavailable.",
		Code:     "SERVICE_UNAVAILABLE",
		HttpCode: 503,
		Original: e,
	}
}

func NewUnknownError(e error) CustomError {
	return CustomError{
		Message:  "Unhandled Error",
		Desc:     "An unknown error occurred.",
		Code:     "UNKNOWN_ERROR",
		HttpCode: 520,
		Original: e,
	}
}
