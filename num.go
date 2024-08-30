// TODO: add `exclusiveMin/Max` rule

package zed

import (
	"errors"
	"github.com/ogen-go/ogen"
)

var _ fieldTrait = (*NumField[float64])(nil)

var ruleMin = defineRule[float64, int64](
	"min",
	func(f float64, lim int64) bool {
		return f >= float64(lim)
	},
	func(a *rule[float64, int64], schema *ogen.Schema) {
		schema.SetMinimum(&a.value)
	},
)

var ruleMax = defineRule[float64, int64](
	"max",
	func(f float64, lim int64) bool {
		return f <= float64(lim)
	},
	func(a *rule[float64, int64], schema *ogen.Schema) {
		schema.SetMaximum(&a.value)
	},
)

type (
	unsigned interface {
		uint8 | uint16 | uint32
	}
	integer interface {
		unsigned | int | int8 | int16 | int32 | int64
	}
	float interface {
		float32 | float64
	}
	NumField[T integer | float] struct {
		fieldTrait
		rules rList[float64]
		err   error
	}
)

func newNumField[T integer | float](mn, mx int64, err string) (f *NumField[T]) {
	f = &NumField[T]{
		rules: make(rList[float64]),
		err:   errors.New(err),
	}
	if mn != 0 || mx != 0 {
		f.rules.add(
			ruleMin(mn, err),
			ruleMax(mx, err),
		)
	}
	return
}

func (f *NumField[T]) Min(val int64, err string) *NumField[T] {
	f.rules.add(ruleMin(val, err))
	return f
}

func (f *NumField[T]) Max(val int64, err string) *NumField[T] {
	f.rules.add(ruleMax(val, err))
	return f
}

func (f *NumField[T]) validate(v any, out any) (e error) {
	val, ok := v.(float64)
	if !ok {
		return f.err
	}
	if e = f.rules.apply(val); e != nil {
		return
	}
	out = T(val)
	return
}

func (f *NumField[T]) toSchema() (s *ogen.Schema) {
	s = ogen.NewSchema()
	switch any(T(0)).(type) {
	case int8:
	case uint8:
	case int16:
	case uint16:
	case int32:
	case int:
		s.SetType("integer").
			SetFormat("int32")
	case uint32:
	case int64:
		s.SetType("integer").
			SetFormat("int64")
	case float32:
		s.SetType("number").
			SetFormat("float")
	case float64:
		s.SetType("number").
			SetFormat("double")
	}
	for _, r := range f.rules {
		r.interceptSchema(s)
	}
	return
}
