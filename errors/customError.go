package errors

import (
	"errors"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrObjectExists = errors.New("object already exists")
	ErrCannotDelete = errors.New("cannot delete")
)
