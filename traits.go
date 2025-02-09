package zed

import "github.com/ogen-go/ogen"

type (
	schemaTrait interface {
		validateGeneric(v any, abortEarly bool) (any, error)
		ToSchema() *ogen.Schema
	}
	ruleTrait[T any] interface {
		name() string
		apply(val T) error
		interceptSchema(*ogen.Schema)
	}
)
