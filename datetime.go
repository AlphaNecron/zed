package zed

import (
	"errors"
	"github.com/ogen-go/ogen"
	"time"
)

var _ Schema[time.Time] = (*DateTimeSchema)(nil)

const (
	EpochNanosecond EpochUnit = 1 << iota
	EpochMicrosecond
	EpochMillisecond
	EpochSecond
)

type (
	DateTimeSchema struct {
		rules     rList[string]
		epochUnit EpochUnit
		layout    string
		err       error
	}
	EpochUnit uint8
)

func newDateTimeSchema(err string) *DateTimeSchema {
	return &DateTimeSchema{
		err:    errors.New(err),
		rules:  make(rList[string]),
		layout: time.RFC3339,
	}
}

func (f *DateTimeSchema) EpochUnit(interval EpochUnit) *DateTimeSchema {
	f.epochUnit = interval
	return f
}

func (f *DateTimeSchema) Layout(layout string) *DateTimeSchema {
	f.layout = layout
	return f
}

func (f *DateTimeSchema) Validate(v any, _ bool) (out time.Time, e error) {
	switch val := v.(type) {
	case string:
		out, e = time.Parse(f.layout, val)
	case float64:
		switch f.epochUnit {
		case EpochNanosecond:
			out = time.Unix(0, int64(val))
		case EpochMicrosecond:
			out = time.UnixMicro(int64(val))
		case EpochMillisecond:
			out = time.UnixMilli(int64(val))
		case EpochSecond:
			out = time.Unix(int64(val), 0)
		default:
			e = f.err
		}
	default:
		e = f.err
	}
	return
}

func (f *DateTimeSchema) validateGeneric(v any, abortEarly bool) (any, error) {
	return f.Validate(v, abortEarly)
}

func (f *DateTimeSchema) ToSchema() *ogen.Schema {
	return ogen.DateTime()
}
