package zed

import (
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
		*baseSchema[string, time.Time]
		epochUnit EpochUnit
		layout    string
	}
	EpochUnit uint8
)

func newDateTimeSchema(err string) *DateTimeSchema {
	return &DateTimeSchema{
		baseSchema: newBaseSchema[string, time.Time](err),
		layout:     time.RFC3339,
	}
}

func (s *DateTimeSchema) EpochUnit(interval EpochUnit) *DateTimeSchema {
	s.epochUnit = interval
	return s
}

func (s *DateTimeSchema) Layout(layout string) *DateTimeSchema {
	s.layout = layout
	return s
}

func (s *DateTimeSchema) Validate(v any, _ SchemaValidationFlag) (out time.Time, e error) {
	switch val := v.(type) {
	case string:
		out, e = time.Parse(s.layout, val)
	case float64:
		switch s.epochUnit {
		case EpochNanosecond:
			out = time.Unix(0, int64(val))
		case EpochMicrosecond:
			out = time.UnixMicro(int64(val))
		case EpochMillisecond:
			out = time.UnixMilli(int64(val))
		case EpochSecond:
			out = time.Unix(int64(val), 0)
		default:
			e = s.err
		}
	case time.Time:
		out = val
	default:
		e = s.err
	}
	return
}

func (s *DateTimeSchema) ValidateGeneric(v any, flags SchemaValidationFlag) (any, error) {
	return s.Validate(v, flags)
}

func (s *DateTimeSchema) ToSchema() *ogen.Schema {
	return ogen.DateTime()
}
