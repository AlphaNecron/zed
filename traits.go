package zed

import "github.com/ogen-go/ogen"

type (
	fieldTrait interface {
		validate(val, out any) error
		toSchema() *ogen.Schema
	}
	ruleTrait[T any] interface {
		name() string
		apply(val T) error
		interceptSchema(*ogen.Schema)
	}
)
