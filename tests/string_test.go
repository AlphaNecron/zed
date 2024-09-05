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
	f := zed.String("expected string value")
	testMulti[string, string](t, f, testData, false, strEq)
	testMulti[any, string](t, f, antitheses, true, nil)
}

func TestStringLen(t *testing.T) {
	testData := []string{
		"foo",
		"bar",
		"baz",
		"quux",
	}
	antitheses := []string{
		"0",
		"01",
		"01234",
		"darcey",
		"necron",
	}
	f := zed.String("expected string").MinLen(3, "minLen").MaxLen(4, "maxLen")
	testMulti[string, string](t, f, testData, false, strEq)
	testMulti[string, string](t, f, antitheses, true, nil)
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
	f := zed.String("expected string").MinLen(3, "minLen").MaxLen(6, "maxLen")
	testMulti[string, string](t, f, testData, false, strEq)
	testMulti[string, string](t, f, antitheses, true, nil)
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
	f := zed.String("expected string").Pattern("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$", "regex")
	testMulti[string, string](t, f, testData, false, strEq)
	testMulti[string, string](t, f, antitheses, true, nil)
}
