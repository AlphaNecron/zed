package zed

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Path string
	Err  error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("`%s`: %v", e.Path, e.Err)
}

func newValidationErr(err error, par string) error {
	if err == nil {
		return nil
	}
	var ve *ValidationError
	if errors.As(err, &ve) {
		ve.Path = par + "." + ve.Path
	} else {
		ve = &ValidationError{
			Path: par,
			Err:  err,
		}
	}
	return ve
}
