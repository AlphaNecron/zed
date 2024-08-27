package tests

import (
	"github.com/stretchr/testify/assert"
	"necron.dev/zed"
	"testing"
)

func TestString(t *testing.T) {
	var foo string
	v := zed.New().
		Field("foo", zed.String())
	assert.Nil(t, v.Validate(
		map[string]any{
			"foo": "foo",
		},
		map[string]any{
			"foo": &foo,
		},
	))
	assert.Equal(t, "foo", foo)
}

func TestStringLen(t *testing.T) {
	var foo, bar string
	assert.Nil(t, zed.New().
		Field("foo", zed.String().MinLen(3, "minLen").MaxLen(3, "maxLen")).
		Field("bar", zed.String().MinLen(4, "minLen2").MaxLen(4, "maxLen2")).
		Validate(
			map[string]any{
				"foo": "bar",
				"bar": "quux",
			},
			map[string]any{
				"foo": &foo,
				"bar": &bar,
			},
		))
	assert.Equal(t, "bar", foo)
	assert.Equal(t, "quux", bar)
}

func TestStringLen2(t *testing.T) {
	testData1 := []string{
		"123",
		"1234",
		"12345",
		"123456",
	}
	testData2 := []string{
		"1",
		"12",
		"1234567",
	}
	v := zed.New().
		Field("foo", zed.String().MinLen(3, "minLen").MaxLen(6, "maxLen"))
	for _, test := range testData1 {
		var foo string
		assert.Nil(t, v.Validate(
			map[string]any{
				"foo": test,
			},
			map[string]any{
				"foo": &foo,
			},
		))
		assert.Equal(t, test, foo)
	}
	for _, test := range testData2 {
		var foo string
		assert.NotNil(t, v.Validate(
			map[string]any{
				"foo": test,
			},
			map[string]any{
				"foo": &foo,
			},
		))
		assert.Equal(t, "", foo)
	}
}

func TestStringPattern(t *testing.T) {
	var testData1 = []string{
		"#000",
		"#FFF",
		"#101010",
		"#ffffff",
		"#FFFFFF",
	}
	var testData2 = []string{
		"#0000",
		"#00FF",
		"#F",
		"0",
		"f",
	}
	v := zed.New().
		Field("foo", zed.String().Pattern("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$", "regex"))
	for _, test := range testData1 {
		var foo string
		assert.Nil(t, v.Validate(
			map[string]any{
				"foo": test,
			},
			map[string]any{
				"foo": &foo,
			},
		))
		assert.Equal(t, test, foo)
	}
	for _, test := range testData2 {
		var foo string
		assert.NotNil(t, v.Validate(
			map[string]any{
				"foo": test,
			},
			map[string]any{
				"foo": &foo,
			},
		))
		assert.Equal(t, "", foo)
	}
}

func TestStringMixed(t *testing.T) {

}
