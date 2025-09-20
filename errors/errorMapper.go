package errors

import (
	"errors"

	"encore.dev/beta/errs"
	"gorm.io/gorm"
)

func MapError(err error) *errs.Error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrorNotFound
	}
	return nil
}
