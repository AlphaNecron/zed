package zed

import "github.com/ogen-go/ogen"

type (
	Field[TOut any] interface {
		Validate(val any) (TOut, error)
		ToSchema() *ogen.Schema
	}
	ruleTrait[T any] interface {
		name() string
		apply(val T) error
		interceptSchema(*ogen.Schema)
	}
)
