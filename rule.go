package zed

import (
	"errors"
	"github.com/ogen-go/ogen"
)

var _ ruleTrait[any] = (*rule[any, any])(nil)

type (
	rule[TSrc, TVal any] struct {
		err              error
		value            TVal
		_name            string
		validate         func(TSrc, TVal) bool
		_interceptSchema func(*rule[TSrc, TVal], *ogen.Schema)
	}
	rList[T any] map[string]ruleTrait[T]
)

func defineRule[TSrc, TVal any](
	name string,
	validateFn func(TSrc, TVal) bool,
	schemaInterceptFn func(*rule[TSrc, TVal], *ogen.Schema),
) func(val TVal, err string) ruleTrait[TSrc] {
	return func(val TVal, err string) ruleTrait[TSrc] {
		return &rule[TSrc, TVal]{
			_name:            name,
			err:              errors.New(err),
			value:            val,
			validate:         validateFn,
			_interceptSchema: schemaInterceptFn,
		}
	}
}

func (a *rule[TSrc, TVal]) name() string {
	return a._name
}

func (a *rule[TSrc, TVal]) apply(v TSrc) (e error) {
	if !a.validate(v, a.value) {
		e = a.err
	}
	return
}

func (a *rule[TSrc, TVal]) interceptSchema(schema *ogen.Schema) {
	a._interceptSchema(a, schema)
}

func (l rList[T]) add(rules ...ruleTrait[T]) {
	for _, r := range rules {
		l[r.name()] = r
	}
}

func (l rList[T]) apply(v T) (e error) {
	for _, a := range l {
		if e = a.apply(v); e != nil {
			return
		}
	}
	return
}
