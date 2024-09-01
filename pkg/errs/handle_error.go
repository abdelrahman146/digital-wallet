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
	// Map gorm errors
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NewNotFoundError("Record not found", err)
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		return NewConflictError("Check constraint violated", err)
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return NewConflictError("Invalid transaction", err)
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return NewBadRequestError("Unsupported relation", err)
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return NewBadRequestError("Primary key required", err)
	case errors.Is(err, gorm.ErrModelValueRequired):
		return NewBadRequestError("Model value required", err)
	case errors.Is(err, gorm.ErrInvalidData):
		return NewBadRequestError("Invalid data", err)
	case errors.Is(err, gorm.ErrRegistered):
		return NewBadRequestError("Registered", err)
	case errors.Is(err, gorm.ErrInvalidField):
		return NewBadRequestError("Invalid field", err)
	case errors.Is(err, gorm.ErrInvalidValue):
		return NewBadRequestError("Invalid value", err)
	case errors.Is(err, gorm.ErrInvalidValueOfLength):
		return NewBadRequestError("Invalid value of length", err)
	case errors.Is(err, gorm.ErrPreloadNotAllowed):
		return NewInternalError("Preload not allowed", err)
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return NewConflictError("Duplicated key", err)
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return NewConflictError("Foreign key violated", err)
	// default
	default:
		return NewUnknownError(err)
	}
}
