package zed

import (
	"github.com/google/uuid"
	"github.com/ogen-go/ogen"
	"regexp"
)

var _ fieldTrait = (*StringField)(nil)

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

var ruleUuid = defineRule[string, any](
	"uuid",
	func(s string, _ any) bool {
		return uuid.Validate(s) == nil
	},
	func(r *rule[string, any], schema *ogen.Schema) {
		schema.SetFormat("uuid")
	},
)

type StringField struct {
	fieldTrait
	rules rList[string]
}

func newStrField() *StringField {
	return &StringField{
		rules: make(rList[string]),
	}
}

func (f *StringField) MinLen(l uint64, err string) *StringField {
	f.rules.add(ruleMinLen(l, err))
	return f
}

func (f *StringField) MaxLen(l uint64, err string) *StringField {
	f.rules.add(ruleMaxLen(l, err))
	return f
}

func (f *StringField) Pattern(p string, err string) *StringField {
	f.rules.add(rulePattern(regexp.MustCompile(p), err))
	return f
}

func (f *StringField) UUID(err string) *StringField {
	f.rules.add(ruleUuid(nil, err))
	return f
}

func (f *StringField) validate(v any, out any) (e error) {
	val := v.(string)
	for _, a := range f.rules {
		if e = a.apply(val); e != nil {
			return
		}
	}
	switch x := out.(type) {
	case *uuid.UUID:
		*x = uuid.MustParse(val)
	case *string:
		*x = val
	}
	return
}

func (f *StringField) toSchema() (s *ogen.Schema) {
	s = ogen.NewSchema()
	for _, r := range f.rules {
		r.interceptSchema(s)
	}
	return
}
