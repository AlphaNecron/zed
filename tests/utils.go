package tests

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"necron.dev/zed"
	"testing"
)

type EqualFunc[TVal, TOut any] func(a TVal, b TOut) (actual, expected TOut, eq bool)

func strEq(a, b string) (string, string, bool) {
	return b, a, a == b
}

func uuidEq(a string, b uuid.UUID) (uuid.UUID, uuid.UUID, bool) {
	return b, uuid.MustParse(a), a == b.String()
}

func boolEq(a, b bool) (bool, bool, bool) {
	return b, a, a == b
}

func testOne[TVal, TOut any](t *testing.T, f zed.Field[TOut], inp TVal, assertErr bool, equalFn EqualFunc[TVal, TOut]) {
	out, e := f.Validate(inp)
	if assertErr {
		assert.Error(t, e)
	} else {
		assert.NoError(t, e)
		if equalFn != nil {
			_actual, expected, eq := equalFn(inp, out)
			assert.Truef(t, eq, "mismatched inp and out, actual=%v, expected=%v", _actual, expected)
		}
	}
}

func testMulti[TVal, TOut any](t *testing.T, f zed.Field[TOut], data []TVal, assertErr bool, equalFn EqualFunc[TVal, TOut]) {
	for _, datum := range data {
		testOne[TVal, TOut](t, f, datum, assertErr, equalFn)
	}
}
