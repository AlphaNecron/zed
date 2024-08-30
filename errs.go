package zed

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	path string
	err  error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("path=%s, error=%s", e.path, e.err)
}

func NewError(e error, path string) error {
	return &ValidationError{
		err:  e,
		path: path,
	}
}

var ErrOutFieldMissing = errors.New("missing out field")
var ErrUnexpectedOutType = errors.New("unexpected type for out")
