package errors

import (
	"fmt"

	"encore.dev/beta/errs"
)

func NewEncoreError(msg string, code errs.ErrCode) *errs.Error {
	return &errs.Error{
		Code:    code,
		Message: msg,
	}
}

var ErrorNotFound = NewEncoreError("not found", errs.NotFound)

var ErrorNil = NewEncoreError("must not be nil", errs.InvalidArgument)

var ErrorIDMissing = NewEncoreError("id must be set", errs.InvalidArgument)

func ErrorAttributeMustBeSet(attribute string) *errs.Error {
	return NewEncoreError(attribute+" must be set", errs.InvalidArgument)
}

var ErrorUnauthenticated = NewEncoreError("unauthenticated", errs.Unauthenticated)

func BadRequest(msg string) *errs.Error {
	return NewEncoreError(msg, errs.InvalidArgument)
}

func BadRequestf(format string, args ...any) *errs.Error {
	msg := fmt.Sprintf(format, args...)
	return NewEncoreError(msg, errs.InvalidArgument)
}

var PiidMissing = NewEncoreError("piid missing", errs.InvalidArgument)

func ErrorReferenceNotFound(err error) *errs.Error {
	return NewEncoreError("reference not found: "+err.Error(), errs.NotFound)
}
