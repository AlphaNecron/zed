// TODO: add `exclusiveMin/Max` rule

package zed

import (
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
		*baseSchema[float64, T]
	}
)

func newNumSchema[T integer | float](mn, mx int64, err string) (f *NumSchema[T]) {
	f = &NumSchema[T]{
		baseSchema: newBaseSchema[float64, T](err),
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

func (s *NumSchema[T]) Min(val int64, exclusive bool, err string) *NumSchema[T] {
	s.rules.add(ruleMin(numLimit{
		value:     val,
		exclusive: exclusive,
	}, err))
	return s
}

func (s *NumSchema[T]) Max(val int64, exclusive bool, err string) *NumSchema[T] {
	s.rules.add(ruleMax(numLimit{
		value:     val,
		exclusive: exclusive,
	}, err))
	return s
}

func (s *NumSchema[T]) Validate(v any, flags SchemaValidationFlag) (out T, e error) {
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
		if e = s.rules.apply(float64(val), flags.Has(AbortEarly)); e == nil {
			out = T(val)
		}
	default:
		e = s.err
	}
	return
}

func (s *NumSchema[T]) validateGeneric(v any, flags SchemaValidationFlag) (any, error) {
	return s.Validate(v, flags)
}

func (s *NumSchema[T]) ToSchema() (os *ogen.Schema) {
	os = ogen.NewSchema()
	switch any(T(0)).(type) {
	case int8:
	case uint8:
	case int16:
	case uint16:
	case int32:
	case int:
		os.SetType("integer").
			SetFormat("int32")
	case uint32:
	case int64:
		os.SetType("integer").
			SetFormat("int64")
	case float32:
		os.SetType("number").
			SetFormat("float")
	case float64:
		os.SetType("number").
			SetFormat("double")
	}
	for _, r := range s.rules {
		r.interceptSchema(os)
	}
	return
}
