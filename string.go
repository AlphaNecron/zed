package zed

import (
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
	*baseSchema[string, string]
}

func newStrSchema(err string) *StringSchema {
	return &StringSchema{
		baseSchema: newBaseSchema[string, string](err),
	}
}

func (s *StringSchema) MinLen(l uint64, err string) *StringSchema {
	s.rules.add(ruleMinLen(l, err))
	return s
}

func (s *StringSchema) MaxLen(l uint64, err string) *StringSchema {
	s.rules.add(ruleMaxLen(l, err))
	return s
}

func (s *StringSchema) Pattern(p string, err string) *StringSchema {
	s.rules.add(rulePattern(regexp.MustCompile(p), err))
	return s
}

func (s *StringSchema) Validate(v any, flags SchemaValidationFlag) (out string, e error) {
	val, vOk := v.(string)
	if !vOk {
		e = s.err
		return
	}
	e = s.rules.apply(val, flags.Has(AbortEarly))
	out = val
	return
}

func (s *StringSchema) ValidateGeneric(v any, flags SchemaValidationFlag) (any, error) {
	return s.Validate(v, flags)
}

func (s *StringSchema) ToSchema() (os *ogen.Schema) {
	os = ogen.String()
	for _, r := range s.rules {
		r.interceptSchema(os)
	}
	return
}
