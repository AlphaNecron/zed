package zed

import (
	"errors"
	"github.com/ogen-go/ogen"
	"regexp"
)

var _ Schema[string] = (*StringSchema)(nil)

var ruleMinLen = defineRule[string, uint64](
	"minLen",
	func(s string, u uint64) bool {
		return len(s) >= int(u)
	},
	func(r *rule[string, uint64], schema *ogen.Schema) {
		schema.SetMinLength(&r.value)
	},
)

var ruleMaxLen = defineRule[string, uint64](
	"maxLen",
	func(s string, u uint64) bool {
		return len(s) <= int(u)
	},
	func(r *rule[string, uint64], schema *ogen.Schema) {
		schema.SetMaxLength(&r.value)
	},
)

var rulePattern = defineRule[string, *regexp.Regexp](
	"pattern",
	func(s string, r *regexp.Regexp) bool {
		return r.MatchString(s)
	},
	func(r *rule[string, *regexp.Regexp], schema *ogen.Schema) {
		schema.SetPattern(r.value.String())
	},
)

type StringSchema struct {
	err   error
	rules rList[string]
}

func newStrSchema(err string) *StringSchema {
	return &StringSchema{
		rules: make(rList[string]),
		err:   errors.New(err),
	}
}

func (f *StringSchema) MinLen(l uint64, err string) *StringSchema {
	f.rules.add(ruleMinLen(l, err))
	return f
}

func (f *StringSchema) MaxLen(l uint64, err string) *StringSchema {
	f.rules.add(ruleMaxLen(l, err))
	return f
}

func (f *StringSchema) Pattern(p string, err string) *StringSchema {
	f.rules.add(rulePattern(regexp.MustCompile(p), err))
	return f
}

func (f *StringSchema) Validate(v any, abortEarly bool) (out string, e error) {
	val, vOk := v.(string)
	if !vOk {
		e = f.err
		return
	}
	e = f.rules.apply(val, abortEarly)
	out = val
	return
}

func (f *StringSchema) validateGeneric(v any, abortEarly bool) (out any, e error) {
	return f.Validate(v, abortEarly)
}

func (f *StringSchema) ToSchema() (s *ogen.Schema) {
	s = ogen.String()
	for _, r := range f.rules {
		r.interceptSchema(s)
	}
	return
}
