package tests

import (
	"maps"
	"necron.dev/zed"
	"slices"
	"testing"
)

func TestBool(t *testing.T) {
	testData := map[any]bool{
		"true":  true,
		true:    true,
		"false": false,
		false:   false,
		"TRUE":  true,
		"FALSE": false,
	}
	actualTestData := slices.Collect(maps.Keys(testData))
	antitheses := []any{
		"foo",
		"bar",
		0,
		1,
	}
	z := zed.New().
		Field("foo", zed.Bool("expected bool value", false))
	testMulti[any, bool](t, z, "foo", actualTestData, false, func(a any, b bool) bool {
		return testData[a] == b
	})
	testMulti[any, bool](t, z, "foo", antitheses, true, nil)
}

func TestBoolStrict(t *testing.T) {
	testData := []bool{
		true,
		false,
	}
	antitheses := []any{
		"foo",
		"bar",
		0,
		1,
		"true",
		"false",
		"TRUE",
		"FALSE",
	}
	z := zed.New().
		Field("foo", zed.Bool("expected bool value", true))
	testMulti[bool, bool](t, z, "foo", testData, false, boolEq)
	testMulti[any, bool](t, z, "foo", antitheses, true, nil)
}
