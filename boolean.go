package zed

import (
	"errors"
	"strings"
)

var _ Field[bool] = (*BoolField)(nil)

type BoolField struct {
	Field[bool]
	rules  rList[string]
	strict bool
	err    error
}

func newBoolField(err string) *BoolField {
	return &BoolField{
		err:   errors.New(err),
		rules: make(rList[string]),
	}
}

func (f *BoolField) Strict() *BoolField {
	f.strict = true
	return f
}

func (f *BoolField) Validate(v any) (out bool, e error) {
	switch val := v.(type) {
	case bool:
		out = val
	case string:
		if f.strict {
			e = f.err
			return
		}
		if strings.EqualFold(val, "true") {
			out = true
		} else if strings.EqualFold(val, "false") {
			out = false
		} else {
			e = f.err
		}
	default:
		e = f.err
	}
	return
}
