package tests

import (
	"necron.dev/zed"
	"testing"
)

func TestString(t *testing.T) {
	testData := []string{
		"foo",
		"bar",
		"baz",
	}
	antitheses := []any{
		true,
		false,
		0,
		1,
	}
	z := zed.New().
		Field("foo", zed.String("expected string value"))
	testMulti[string, string](t, z, "foo", testData, false, strEq)
	testMulti[any, string](t, z, "foo", antitheses, true, nil)
}

func TestStringLen(t *testing.T) {
	testOne[string, string](t, zed.New().
		Field("foo", zed.String("expected string").MinLen(3, "minLen").MaxLen(3, "maxLen")).
		Field("bar", zed.String("expected string").MinLen(4, "minLen2").MaxLen(4, "maxLen2")),
		map[string]any{
			"foo": "bar",
			"bar": "quux",
		},
		false, strEq)
}

func TestStringLen2(t *testing.T) {
	testData := []string{
		"123",
		"1234",
		"12345",
		"123456",
	}
	antitheses := []string{
		"1",
		"12",
		"1234567",
	}
	z := zed.New().
		Field("foo", zed.String("expected string").MinLen(3, "minLen").MaxLen(6, "maxLen"))
	testMulti[string, string](t, z, "foo", testData, false, strEq)
	testMulti[string, string](t, z, "foo", antitheses, true, nil)
}

func TestStringPattern(t *testing.T) {
	testData := []string{
		"#000",
		"#FFF",
		"#101010",
		"#ffffff",
		"#FFFFFF",
	}
	antitheses := []string{
		"#0000",
		"#00FF",
		"#F",
		"0",
		"f",
	}
	z := zed.New().
		Field("foo", zed.String("expected string").Pattern("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$", "regex"))
	testMulti[string, string](t, z, "foo", testData, false, strEq)
	testMulti[string, string](t, z, "foo", antitheses, true, nil)
}
