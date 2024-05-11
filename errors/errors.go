package errors

import "encore.dev/beta/errs"

func NewError(msg string, code errs.ErrCode) error {
	return &errs.Error{
		Code:    code,
		Message: msg,
	}
}

var ErrorNotFound = NewError("not found", errs.NotFound)

var ErrorNil = NewError("must not be nil", errs.InvalidArgument)

var ErrorIDMissing = NewError("id must be set", errs.InvalidArgument)

func ErrorAttributeMustBeSet(attribute string) error {
	return NewError(attribute+" must be set", errs.InvalidArgument)
}

var ErrorUnauthenticated = NewError("unauthenticated", errs.Unauthenticated)
