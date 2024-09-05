package zed

import "github.com/ogen-go/ogen"

var _ Field[any] = (*StructField[any])(nil)

type StructField[T any] struct {
	Field[T]
}

func (f *StructField[T]) Validate(m any) (out T, e error) {
	return
}

func (f *StructField[T]) ToSchema() (s *ogen.Schema) {
	return
}
