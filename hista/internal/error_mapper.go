package internal

import (
	"errors"

	"encore.dev/beta/errs"
	"gorm.io/gorm"
)

func errorMapper(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return &errs.Error{Code: errs.NotFound, Message: "not found"}
	default:
		return &errs.Error{Code: errs.Internal, Message: "internal database error", Details: errs.Details(err)}
	}
}
