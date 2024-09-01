package zed

import (
	"errors"
	"strings"
)

var _ fieldTrait = (*BoolField)(nil)

type BoolField struct {
	fieldTrait
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

func (f *BoolField) validate(v any, out any) (e error) {
	o, ok := out.(*bool)
	if !ok {
		return ErrUnexpectedOutType
	}
	switch val := v.(type) {
	case bool:
		*o = val
	case string:
		if f.strict {
			return f.err
		}
		if strings.EqualFold(val, "true") {
			*o = true
		} else if strings.EqualFold(val, "false") {
			*o = false
		} else {
			e = f.err
		}
	default:
		return f.err
	}
	return
}
