package tests

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"necron.dev/zed"
	"slices"
	"testing"
)

type AuthReq struct {
	Username string `zed:"username,err:'username is required',pattern:'^[a-zA-Z0-9_]{3,20}$',pattern_err:'invalid username'"`
	Password string `zed:"password,err:'password is required',min_len:6,max_len:20,min_len_err:'password is too short',max_len_err:'password is too long'"`
	Email    string `zed:"email,err:'email is required',pattern:'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$',pattern_err:'invalid email'"`
	PIN      string `zed:"pin,err:'pin is required',min_len:4,max_len:4,min_len_err:'pin must be 4 digits',max_len_err:'pin must be 4 digits'"`
}

func TestStructConstructionShallow(t *testing.T) {
	_, e := zed.StructForE[AuthReq]("expected struct")
	assert.NoError(t, e)
}

func TestStructShallow(t *testing.T) {
	f, e := zed.StructForE[AuthReq]("expected struct")
	assert.NoError(t, e)
	testData := []any{
		AuthReq{
			Username: "foo",
			Password: "bar123_",
			Email:    "baz@foo.bar",
			PIN:      "1234",
		},
	}
	antitheses := []any{
		"not a struct",
		1,
		false,
		AuthReq{
			Username: "foo@@",
			Password: "bar",
			Email:    "baz",
			PIN:      "123",
		},
		AuthReq{
			Username: "foo",
			Password: "barfoo123",
		},
		AuthReq{
			Username: "foo",
			Password: "bar",
			PIN:      "3456",
		},
	}
	testMultiStruct[AuthReq](t, f, testData, nil)
	testMultiStruct[AuthReq](t, f, antitheses, [][]string{
		{""},
		{""},
		{""},
		{"username", "password", "email", "pin"},
		{"email", "pin"},
		{"password", "email"},
	})
}

func testMultiStruct[TOut any](t *testing.T, f *zed.StructSchema[TOut], testData []any, assertErrPaths [][]string) {
	for i, test := range testData {
		_, e := f.Validate(test, false)
		if len(assertErrPaths) == 0 {
			assert.NoError(t, e)
		} else {
			if e == nil {
				assert.Error(t, e)
			} else {
				for _, _e := range multierr.Errors(e) {
					var validationError *zed.ValidationError
					if errors.As(_e, &validationError) {
						path := validationError.Path
						if !slices.Contains(assertErrPaths[i], path) {
							assert.NoError(t, _e)
						} else {
							assertErrPaths[i] = slices.DeleteFunc(assertErrPaths[i], func(s string) bool {
								return s == path
							})
						}
					}
				}
			}
		}
	}
}
