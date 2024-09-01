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
		return NewNotFoundError(err.Error(), nil)
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		return NewConflictError(err.Error(), nil)
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return NewConflictError(err.Error(), nil)
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		return NewBadRequestError(err.Error(), nil)
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return NewBadRequestError(err.Error(), nil)
	case errors.Is(err, gorm.ErrModelValueRequired):
		return NewBadRequestError(err.Error(), nil)
	case errors.Is(err, gorm.ErrInvalidData):
		return NewBadRequestError(err.Error(), nil)
	case errors.Is(err, gorm.ErrRegistered):
		return NewBadRequestError(err.Error(), nil)
	case errors.Is(err, gorm.ErrInvalidField):
		return NewBadRequestError(err.Error(), nil)
	case errors.Is(err, gorm.ErrInvalidValue):
		return NewBadRequestError(err.Error(), nil)
	case errors.Is(err, gorm.ErrInvalidValueOfLength):
		return NewBadRequestError(err.Error(), nil)
	case errors.Is(err, gorm.ErrPreloadNotAllowed):
		return NewInternalError(err.Error(), nil)
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return NewConflictError(err.Error(), nil)
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return NewConflictError(err.Error(), nil)
	// default
	default:
		return NewUnknownError(err)
	}
}
