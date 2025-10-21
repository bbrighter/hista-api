package errors

import (
	"errors"
	"fmt"
)

type CustomError struct {
	Kind    error
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s, kind=%s", e.Message, e.Kind.Error())
}

var (
	ErrBadRequest = errors.New("bad request")
)

func NewCustomError(message string, err error) *CustomError {
	return &CustomError{Message: message, Kind: err}
}

func NewErrBadRequest(message string) *CustomError {
	return NewCustomError(message, ErrBadRequest)
}
