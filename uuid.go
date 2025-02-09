package zed

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen"
)

var _ Schema[uuid.UUID] = (*UUIDSchema)(nil)

type UUIDSchema struct {
	rules rList[string]
	err   error
}

func newUuidSchema(err string) *UUIDSchema {
	return &UUIDSchema{
		err:   errors.New(err),
		rules: make(rList[string]),
	}
}

func (f *UUIDSchema) Validate(v any, _ bool) (out uuid.UUID, e error) {
	val, vOk := v.(string)
	if !vOk {
		e = f.err
		return
	}
	out, e = uuid.Parse(val)
	if e != nil {
		e = f.err
	}
	return
}

func (f *UUIDSchema) validateGeneric(v any, abortEarly bool) (out any, e error) {
	return f.Validate(v, abortEarly)
}

func (f *UUIDSchema) ToSchema() *ogen.Schema {
	return ogen.UUID()
}
