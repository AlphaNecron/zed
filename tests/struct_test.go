package tests

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
	"necron.dev/zed"
	"slices"
	"testing"
	"time"
)

type (
	AuthReq struct {
		Username        string `zed:"username,err:'invalid username',pattern:'^[a-zA-Z0-9_]{3,20}$',pattern_err:'invalid username',required"`
		Password        string `zed:"password,err:'invalid password',min_len:6,max_len:20,min_len_err:'password is too short',max_len_err:'password is too long',required"`
		Email           string `zed:"email,err:'invalid email',pattern:'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$',pattern_err:'invalid email',required"`
		RememberSession bool   `zed:"rememberSession,err:'invalid rememberSession'"`
	}
	NotSoComplexObject struct {
		Username        string    `zed:"username,err:'invalid username',pattern:'^[a-zA-Z0-9_]{6,20}$',pattern_err:'invalid username',min_len:6,max_len:20,min_len_err:'username is too short',max_len_err:'username is too long',required"`
		Password        string    `zed:"password,err:'invalid password',min_len:6,max_len:20,min_len_err:'password is too short',max_len_err:'password is too long',required"`
		Email           string    `zed:"email,err:'invalid email',pattern:'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$',pattern_err:'invalid email',required"`
		RememberSession bool      `zed:"rememberSession,err:'invalid rememberSession'"`
		Age             int       `zed:"age,err:'invalid age',min:0,exclusive_min,max:100,min_err:'invalid age',max_err:'invalid age'"`
		Birthdate       time.Time `zed:"birthdate,err:'invalid birthdate',required"`
		UUID            uuid.UUID `zed:"uuid,err:'invalid uuid',required"`
	}
)

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
		},
		AuthReq{
			Username:        "quux",
			Password:        "baz@456",
			Email:           "foo@bar.qux",
			RememberSession: true,
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
		},
		AuthReq{
			Username: "foo",
			Password: "barfoo123",
		},
		AuthReq{
			Username: "foo",
			Password: "bar",
		},
	}
	testMultiStruct[AuthReq](t, f, testData, nil)
	testMultiStruct[AuthReq](t, f, antitheses, [][]string{
		{""},
		{""},
		{""},
		{"username", "password", "email"},
		{"email"},
		{"password", "email"},
	})
}

func TestStructFromMap(t *testing.T) {
	f, se := zed.StructForE[AuthReq]("expected struct")
	assert.NoError(t, se)
	testData := []any{
		map[string]any{
			"username": "foo",
			"password": "bar123_",
			"email":    "baz@foo.bar",
		},
		map[string]any{
			"username":              "foo",
			"password":              "bar123_",
			"email":                 "baz@foo.bar",
			"an arbitrary key":      time.Now(),
			"another arbitrary key": uuid.New(),
		},
	}
	expected := AuthReq{
		Username: "foo",
		Password: "bar123_",
		Email:    "baz@foo.bar",
	}
	for _, test := range testData {
		out, e := f.Validate(test, zed.AbortEarly)
		assert.NoError(t, e)
		assert.Equal(t, expected, out)
	}
}

func TestStructStripMap(t *testing.T) {
	f := zed.Struct("expected map").
		AddField("username", zed.String("invalid username"), zed.FieldRequired).
		AddField("password", zed.String("invalid password"), zed.FieldRequired).
		AddField("email", zed.String("invalid email"), zed.FieldRequired)
	testData := []any{
		map[string]any{
			"username": "foo",
			"password": "bar123_",
			"email":    "baz@foo.bar",
		},
		map[string]any{
			"username":              "foo",
			"password":              "bar123_",
			"email":                 "baz@foo.bar",
			"an arbitrary key":      time.Now(),
			"another arbitrary key": uuid.New(),
		},
	}
	expected := map[string]any{
		"username": "foo",
		"password": "bar123_",
		"email":    "baz@foo.bar",
	}
	for _, test := range testData {
		out, e := f.Validate(test, zed.AbortEarly)
		assert.NoError(t, e)
		assert.Equal(t, expected, out)
	}
}

func BenchmarkStructCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, e := zed.StructForE[NotSoComplexObject]("expected struct")
		assert.NoError(b, e)
	}
}

func BenchmarkStructValidation(b *testing.B) {
	f, e := zed.StructForE[NotSoComplexObject]("expected struct")
	assert.NoError(b, e)
	now := time.Now()
	mockUuid := uuid.New().String()
	for i := 0; i < b.N; i++ {
		_, e := f.Validate(map[string]any{
			"username":        "foobar",
			"password":        "bar123_",
			"email":           "baz@foo.bar",
			"age":             18,
			"birthdate":       now,
			"rememberSession": true,
			"uuid":            mockUuid,
		}, zed.AbortEarly)
		assert.NoError(b, e)
	}
}

func testMultiStruct[TOut any](t *testing.T, f *zed.StructSchema[TOut], testData []any, assertErrPaths [][]string) {
	for i, test := range testData {
		_, e := f.Validate(test, 0)
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
