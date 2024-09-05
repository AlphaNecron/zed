package zed

import (
	"errors"
	"github.com/ogen-go/ogen"
	"regexp"
)

var _ Field[string] = (*StringField)(nil)

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

type StringField struct {
	Field[string]
	err   error
	rules rList[string]
}

func newStrField(err string) *StringField {
	return &StringField{
		rules: make(rList[string]),
		err:   errors.New(err),
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

func (f *StringField) Validate(v any) (out string, e error) {
	val, vOk := v.(string)
	if !vOk {
		e = f.err
		return
	}
	for _, a := range f.rules {
		if e = a.apply(val); e != nil {
			return
		}
	}
	out = val
	return
}

func (f *StringField) ToSchema() (s *ogen.Schema) {
	s = ogen.NewSchema().SetType("string")
	for _, r := range f.rules {
		r.interceptSchema(s)
	}
	return
}
