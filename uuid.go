package zed

import (
	"errors"
	"github.com/google/uuid"
)

var _ Field[uuid.UUID] = (*UUIDField)(nil)

type UUIDField struct {
	Field[uuid.UUID]
	rules rList[string]
	err   error
}

func newUuidField(err string) *UUIDField {
	return &UUIDField{
		err:   errors.New(err),
		rules: make(rList[string]),
	}
}

func (f *UUIDField) Validate(v any) (out uuid.UUID, e error) {
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
