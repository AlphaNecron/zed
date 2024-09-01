package zed

import (
	"errors"
	"time"
)

var _ fieldTrait = (*DateTimeField)(nil)

const (
	EpochNanosecond EpochUnit = 1 << iota
	EpochMicrosecond
	EpochMillisecond
	EpochSecond
)

type (
	DateTimeField struct {
		fieldTrait
		rules     rList[string]
		epochUnit EpochUnit
		layout    string
		err       error
	}
	EpochUnit uint8
)

func newDateTimeField(err string) *DateTimeField {
	return &DateTimeField{
		err:    errors.New(err),
		rules:  make(rList[string]),
		layout: time.RFC3339,
	}
}

func (f *DateTimeField) EpochUnit(interval EpochUnit) *DateTimeField {
	f.epochUnit = interval
	return f
}

func (f *DateTimeField) Layout(layout string) *DateTimeField {
	f.layout = layout
	return f
}

func (f *DateTimeField) validate(v any, out any) (e error) {
	o, ok := out.(*time.Time)
	if !ok {
		return ErrUnexpectedOutType
	}
	switch val := v.(type) {
	case string:
		*o, e = time.Parse(f.layout, val)
	case float64:
		switch f.epochUnit {
		case EpochNanosecond:
			*o = time.Unix(0, int64(val))
			break
		case EpochMicrosecond:
			*o = time.UnixMicro(int64(val))
			break
		case EpochMillisecond:
			*o = time.UnixMilli(int64(val))
			break
		case EpochSecond:
			*o = time.Unix(int64(val), 0)
			break
		default:
			e = f.err
		}
	default:
		e = f.err
	}
	return
}
