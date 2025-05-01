package zed

import (
	"github.com/google/uuid"
	"github.com/ogen-go/ogen"
)

var _ Schema[uuid.UUID] = (*UUIDSchema)(nil)

type UUIDSchema struct {
	*baseSchema[any, uuid.UUID]
}

func newUuidSchema(err string) *UUIDSchema {
	return &UUIDSchema{
		baseSchema: newBaseSchema[any, uuid.UUID](err),
	}
}

func (s *UUIDSchema) Validate(v any, _ SchemaValidationFlag) (out uuid.UUID, e error) {
	switch val := v.(type) {
	case string:
		out, e = uuid.Parse(val)
	case uuid.UUID:
		out = val
	default:
		e = s.err
	}
	return
}

func (s *UUIDSchema) ValidateGeneric(v any, flags SchemaValidationFlag) (any, error) {
	return s.Validate(v, flags)
}

func (s *UUIDSchema) ToSchema() *ogen.Schema {
	return ogen.UUID()
}
