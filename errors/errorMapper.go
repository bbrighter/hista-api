package errors

import (
	"errors"

	"encore.dev/beta/errs"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func MapError(err error) error {
	if err == nil {
		print("error is nil")
		return nil
	}

	print("error is not nil")
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrorNotFound
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return ErrorUnauthenticated
	default:
		return &errs.Error{Code: errs.Internal, Message: "internal error", Details: errs.Details(err)}
	}
}
