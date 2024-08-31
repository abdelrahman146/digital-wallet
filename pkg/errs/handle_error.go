package errs

import "errors"

func HandleError(err error) CustomError {
	var customError CustomError
	switch {
	case errors.As(err, &customError):
		return err.(CustomError)
	default:
		return NewUnknownError(err)
	}
}
