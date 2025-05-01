package zed

import "github.com/ogen-go/ogen"

type (
	SchemaTrait interface {
		ValidateGeneric(v any, flags SchemaValidationFlag) (any, error)
		ToSchema() *ogen.Schema
	}
	ruleTrait[T any] interface {
		name() string
		apply(val T) error
		interceptSchema(*ogen.Schema)
	}
)
