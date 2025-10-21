package errors

import (
	"errors"

	"encore.dev/beta/errs"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func MapError(err error) error {

	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrorNotFound
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return ErrorUnauthenticated
	case errors.Is(err, PiidMissing):
		return PiidMissing
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return NewEncoreError("duplicate key", errs.AlreadyExists)

	}

	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		switch pgErr.Code {
		case "25505":
			return NewEncoreError(pgErr.Error(), errs.AlreadyExists)
		case "23505":
			return NewEncoreError(pgErr.Error(), errs.AlreadyExists)
		case "23503":
			return NewEncoreError(pgErr.Error(), errs.NotFound)
		}
	}

	customerErr, ok := err.(*CustomError)
	if ok {
		switch customerErr.Kind {
		case ErrBadRequest:
			return NewEncoreError(customerErr.Error(), errs.InvalidArgument)
		}
	}

	encoreErr, ok := err.(*errs.Error)
	if ok {
		return encoreErr
	}

	return &errs.Error{Code: errs.Internal, Message: "internal error", Details: errs.Details(err)}

}
