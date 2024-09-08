package errors

import "encore.dev/beta/errs"

func NewError(msg string, code errs.ErrCode) *errs.Error {
	return &errs.Error{
		Code:    code,
		Message: msg,
	}
}

var ErrorNotFound = NewError("not found", errs.NotFound)

var ErrorNil = NewError("must not be nil", errs.InvalidArgument)

var ErrorIDMissing = NewError("id must be set", errs.InvalidArgument)

func ErrorAttributeMustBeSet(attribute string) *errs.Error {
	return NewError(attribute+" must be set", errs.InvalidArgument)
}

var ErrorUnauthenticated = NewError("unauthenticated", errs.Unauthenticated)

func BadRequest(msg string) *errs.Error {
	return NewError(msg, errs.InvalidArgument)
}
