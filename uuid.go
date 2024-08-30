package zed

import (
	"errors"
	"github.com/google/uuid"
)

var _ fieldTrait = (*UUIDField)(nil)

type UUIDField struct {
	fieldTrait
	rules rList[string]
	err   error
}

func newUuidField(err string) *UUIDField {
	return &UUIDField{
		err:   errors.New(err),
		rules: make(rList[string]),
	}
}

func (f *UUIDField) validate(v any, out any) (e error) {
	o, outOk := out.(*uuid.UUID)
	if !outOk {
		return ErrUnexpectedOutType
	}
	val, vOk := v.(string)
	if !vOk {
		return f.err
	}
	*o, e = uuid.Parse(val)
	if e != nil {
		return f.err
	}
	return
}
