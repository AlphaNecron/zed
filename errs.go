package zed

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Path string
	Err  error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("path=%s, error=%s", e.Path, e.Err)
}

func NewError(e error, path string) error {
	return &ValidationError{
		Err:  e,
		Path: path,
	}
}

var ErrOutFieldMissing = errors.New("missing out field")
var ErrUnexpectedOutType = errors.New("unexpected type for out")
