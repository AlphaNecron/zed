package tests

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"necron.dev/zed"
	"testing"
)

type EqualFunc[TVal, TOut any] func(a TVal, b TOut) bool

func strEq(a, b string) bool {
	return a == b
}

func uuidEq(a string, b uuid.UUID) bool {
	return a == b.String()
}

func boolEq(a, b bool) bool {
	return a == b
}

func mapPtr[TKey comparable, TVal any](m map[TKey]TVal, k TKey) *TVal {
	v, _ := m[k]
	return &v
}

func takeVal[T any](x any) T {
	return *(x.(*T))
}

func testOne[TVal, TOut any](t *testing.T, z *zed.Zed, inp map[string]any, assertErr bool, equalFn EqualFunc[TVal, TOut]) {
	out := make(map[string]any)
	actual := make(map[string]TOut)
	for k := range inp {
		var _default TOut
		actual[k] = _default
		out[k] = mapPtr(actual, k)
	}
	e := z.Validate(
		inp,
		out,
	)
	if assertErr {
		assert.NotNil(t, e)
	} else {
		assert.Nil(t, e)
		if equalFn != nil {
			for k, v := range inp {
				assert.True(t, equalFn(v.(TVal), takeVal[TOut](out[k])), "mismatched inp and out")
			}
		}
	}
}

func testMulti[TVal, TOut any](t *testing.T, z *zed.Zed, keyName string, data []TVal, assertErr bool, equalFn EqualFunc[TVal, TOut]) {
	for _, datum := range data {
		m := make(map[string]any)
		m[keyName] = datum
		testOne[TVal, TOut](t, z, m, assertErr, equalFn)
	}
}
