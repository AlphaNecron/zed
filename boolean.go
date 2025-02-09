package zed

import (
	"errors"
	"github.com/ogen-go/ogen"
)

var _ Schema[bool] = (*BoolSchema)(nil)

type BoolSchema struct {
	rules rList[string]
	// strict bool
	err error
}

// func (f *BoolSchema) Strict() *BoolSchema {
// 	f.strict = true
// 	return f
// }

func (f *BoolSchema) Validate(v any, _ bool) (out bool, e error) {
	switch val := v.(type) {
	case bool:
		out = val
	// case string:
	// 	if f.strict {
	// 		e = f.Err
	// 		return
	// 	}
	// 	if strings.EqualFold(val, "true") {
	// 		out = true
	// 	} else if strings.EqualFold(val, "false") {
	// 		out = false
	// 	} else {
	// 		e = f.Err
	// 	}
	default:
		e = f.err
	}
	return
}

func (f *BoolSchema) validateGeneric(v any, abortEarly bool) (any, error) {
	return f.Validate(v, abortEarly)
}

func (f *BoolSchema) ToSchema() *ogen.Schema {
	return ogen.Bool()
}

func newBoolSchema(err string) *BoolSchema {
	return &BoolSchema{
		err:   errors.New(err),
		rules: make(rList[string]),
	}
}
