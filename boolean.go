package zed

import (
	"github.com/ogen-go/ogen"
)

var _ Schema[bool] = (*BoolSchema)(nil)

type BoolSchema struct {
	*baseSchema[any, bool]
}

// func (f *BoolSchema) Strict() *BoolSchema {
// 	f.strict = true
// 	return f
// }

func (s *BoolSchema) Validate(v any, _ SchemaValidationFlag) (out bool, e error) {
	switch val := v.(type) {
	case bool:
		out = val
	// case string:
	// 	if s.strict {
	// 		e = s.Err
	// 		return
	// 	}
	// 	if strings.EqualFold(val, "true") {
	// 		out = true
	// 	} else if strings.EqualFold(val, "false") {
	// 		out = false
	// 	} else {
	// 		e = s.Err
	// 	}
	default:
		e = s.err
	}
	return
}

func (s *BoolSchema) ValidateGeneric(v any, flags SchemaValidationFlag) (any, error) {
	return s.Validate(v, flags)
}

func (s *BoolSchema) ToSchema() *ogen.Schema {
	return ogen.Bool()
}

func newBoolSchema(err string) *BoolSchema {
	return &BoolSchema{
		baseSchema: newBaseSchema[any, bool](err),
	}
}
