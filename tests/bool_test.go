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
	f := zed.Bool("expected bool value")
	testMulti[any, bool](t, f, actualTestData, false, func(a any, b bool) (bool, bool, bool) {
		return b, testData[a], testData[a] == b
	})
	testMulti[any, bool](t, f, antitheses, true, nil)
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
	f := zed.Bool("expected bool value").Strict()
	testMulti[bool, bool](t, f, testData, false, boolEq)
	testMulti[any, bool](t, f, antitheses, true, nil)
}
