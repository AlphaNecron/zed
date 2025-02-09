// TODO: add `exclusiveMin/Max` rule

package zed

import (
	"errors"
	"github.com/ogen-go/ogen"
)

type numLimit struct {
	value     int64
	exclusive bool
}

var _ Schema[float64] = (*NumSchema[float64])(nil)

var ruleMin = defineRule[float64, numLimit](
	"min",
	func(f float64, lim numLimit) bool {
		return f > float64(lim.value) || (!lim.exclusive && f == float64(lim.value))
	},
	func(a *rule[float64, numLimit], schema *ogen.Schema) {
		schema.SetMinimum(&a.value.value).SetExclusiveMinimum(a.value.exclusive)
	},
)

var ruleMax = defineRule[float64, numLimit](
	"max",
	func(f float64, lim numLimit) bool {
		return f < float64(lim.value) || (lim.exclusive && f == float64(lim.value))
	},
	func(a *rule[float64, numLimit], schema *ogen.Schema) {
		schema.SetMaximum(&a.value.value).SetExclusiveMaximum(a.value.exclusive)
	},
)

type (
	unsigned interface {
		uint | uint8 | uint16 | uint32
	}
	integer interface {
		unsigned | int | int8 | int16 | int32 | int64
	}
	float interface {
		float32 | float64
	}
	NumSchema[T integer | float] struct {
		rules rList[float64]
		err   error
	}
)

func newNumSchema[T integer | float](mn, mx int64, err string) (f *NumSchema[T]) {
	f = &NumSchema[T]{
		rules: make(rList[float64]),
		err:   errors.New(err),
	}
	if mn != 0 || mx != 0 {
		f.rules.add(
			ruleMin(numLimit{
				value:     mn,
				exclusive: false,
			}, err),
			ruleMax(numLimit{
				value:     mx,
				exclusive: false,
			}, err),
		)
	}
	return
}

func (f *NumSchema[T]) Min(val int64, exclusive bool, err string) *NumSchema[T] {
	f.rules.add(ruleMin(numLimit{
		value:     val,
		exclusive: exclusive,
	}, err))
	return f
}

func (f *NumSchema[T]) Max(val int64, exclusive bool, err string) *NumSchema[T] {
	f.rules.add(ruleMax(numLimit{
		value:     val,
		exclusive: exclusive,
	}, err))
	return f
}

func (f *NumSchema[T]) Validate(v any, abortEarly bool) (out T, e error) {
	switch val := v.(type) {
	case float32:
	case float64:
	case uint8:
	case int8:
	case uint16:
	case int16:
	case uint32:
	case int32:
	case uint64:
	case int64:
	case uint:
	case int:
		if e = f.rules.apply(float64(val), abortEarly); e == nil {
			out = T(val)
		}
	default:
		e = f.err
	}
	return
}

func (f *NumSchema[T]) validateGeneric(v any, abortEarly bool) (out any, e error) {
	return f.Validate(v, abortEarly)
}

func (f *NumSchema[T]) ToSchema() (s *ogen.Schema) {
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
