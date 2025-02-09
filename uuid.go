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
	val, vOk := v.(string)
	if !vOk {
		e = s.err
		return
	}
	out, e = uuid.Parse(val)
	if e != nil {
		e = s.err
	}
	return
}

func (s *UUIDSchema) validateGeneric(v any, flags SchemaValidationFlag) (any, error) {
	return s.Validate(v, flags)
}

func (s *UUIDSchema) ToSchema() *ogen.Schema {
	return ogen.UUID()
}
