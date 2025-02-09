package zed

import (
	"errors"
)

type (
	baseSchema[TRule, TOut any] struct {
		Schema[TOut]
		err   error
		rules ruleset[TRule]
	}
)

func newBaseSchema[TRule any, TOut any](err string) *baseSchema[TRule, TOut] {
	return &baseSchema[TRule, TOut]{
		err:   errors.New(err),
		rules: make(ruleset[TRule]),
	}
}
