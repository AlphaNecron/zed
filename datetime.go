package zed

import (
	"errors"
	"time"
)

var _ Field[time.Time] = (*DateTimeField)(nil)

const (
	EpochNanosecond EpochUnit = 1 << iota
	EpochMicrosecond
	EpochMillisecond
	EpochSecond
)

type (
	DateTimeField struct {
		Field[time.Time]
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

func (f *DateTimeField) Validate(v any) (out time.Time, e error) {
	switch val := v.(type) {
	case string:
		out, e = time.Parse(f.layout, val)
	case float64:
		switch f.epochUnit {
		case EpochNanosecond:
			out = time.Unix(0, int64(val))
			break
		case EpochMicrosecond:
			out = time.UnixMicro(int64(val))
			break
		case EpochMillisecond:
			out = time.UnixMilli(int64(val))
			break
		case EpochSecond:
			out = time.Unix(int64(val), 0)
			break
		default:
			e = f.err
		}
	default:
		e = f.err
	}
	return
}
