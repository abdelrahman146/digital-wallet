package errs

import (
	"errors"
	"gorm.io/gorm"
)

func HandleError(err error) CustomError {
	var customError CustomError
	switch {
	case errors.As(err, &customError):
		return err.(CustomError)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NewNotFoundError(err.Error(), "RECORD_NOT_FOUND", err)
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		return NewConflictError(err.Error(), "CHECK_CONSTRAINT_VIOLATED", err)
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return NewExpectationFailedError(err.Error(), "INVALID_TRANSACTION", err)
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return NewExpectationFailedError(err.Error(), "UNSUPPORTED_RELATION", err)
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return NewBadRequestError(err.Error(), "PRIMARY_KEY_REQUIRED", err)
	case errors.Is(err, gorm.ErrModelValueRequired):
		return NewExpectationFailedError(err.Error(), "MODEL_VALUE_REQUIRED", err)
	case errors.Is(err, gorm.ErrInvalidData):
		return NewBadRequestError(err.Error(), "INVALID_DATA", err)
	case errors.Is(err, gorm.ErrRegistered):
		return NewNotAcceptableError(err.Error(), "REGISTERED", err)
	case errors.Is(err, gorm.ErrInvalidField):
		return NewBadRequestError(err.Error(), "INVALID_FIELD", err)
	case errors.Is(err, gorm.ErrInvalidValue):
		return NewBadRequestError(err.Error(), "INVALID_VALUE", err)
	case errors.Is(err, gorm.ErrInvalidValueOfLength):
		return NewBadRequestError(err.Error(), "INVALID_VALUE_OF_LENGTH", err)
	case errors.Is(err, gorm.ErrPreloadNotAllowed):
		return NewNotAcceptableError(err.Error(), "PRELOAD_NOT_ALLOWED", err)
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return NewConflictError(err.Error(), "DUPLICATED_KEY", err)
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return NewBadRequestError(err.Error(), "FOREIGN_KEY_VIOLATED", err)
	// default
	default:
		return NewUnknownError(err)
	}
}
