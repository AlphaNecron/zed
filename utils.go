package zed

import (
	"errors"
	"github.com/vmihailenco/tagparser/v2"
	"reflect"
	"strconv"
	"strings"
)

type pair[TF, TS any] struct {
	first  TF
	second TS
}

func makePair[TF, TS any](first TF, second TS) pair[TF, TS] {
	return pair[TF, TS]{
		first:  first,
		second: second,
	}
}

func parseStructFieldName(f reflect.StructField) (string, bool) {
	name := f.Name
	t := f.Tag.Get("zed")
	tagName := t[:strings.Index(t, ",")]
	if tagName == "-" {
		return "", false
	}
	if tagName != "" {
		name = tagName
	}
	return name, true
}

func isUuidType(t reflect.Type) bool {
	return t.PkgPath() == "github.com/google/uuid" && t.Name() == "UUID"
}

func isTimeType(t reflect.Type) bool {
	return t.PkgPath() == "time" && t.Name() == "Time"
}

func parseNumFieldOpts[T integer | float](f *NumSchema[T], t *tagparser.Tag) (*NumSchema[T], error) {
	if t.HasOption("min") {
		if !t.HasOption("min_err") {
			return f, errors.New("`min` option is set but `min_err` option is missing")
		}
		v, e := strconv.ParseInt(t.Options["min"], 10, 64)
		if e != nil {
			return f, errors.New("invalid value for `min` option")
		}
		f.Min(v, t.HasOption("exclusive_min"), t.Options["min_err"])
	}
	if t.HasOption("max") {
		if !t.HasOption("max_err") {
			return f, errors.New("`max` option is set but `max_err` option is missing")
		}
		v, e := strconv.ParseInt(t.Options["max"], 10, 64)
		if e != nil {
			return f, errors.New("invalid value for `max` option")
		}
		f.Max(v, t.HasOption("exclusive_max"), t.Options["max_err"])
	}
	return f, nil
}

func parseStringFieldOpts(f *StringSchema, t *tagparser.Tag) (*StringSchema, error) {
	if t.HasOption("min_len") {
		if !t.HasOption("min_len_err") {
			return f, errors.New("`min_len` option is set but `min_len_err` option is missing")
		}
		v, e := strconv.ParseUint(t.Options["min_len"], 10, 64)
		if e != nil {
			return f, errors.New("invalid value for `min_len` option")
		}
		f.MinLen(v, t.Options["min_len_err"])
	}
	if t.HasOption("max_len") {
		if !t.HasOption("max_len_err") {
			return f, errors.New("`max_len` option is set but `max_len_err` option is missing")
		}
		v, e := strconv.ParseUint(t.Options["max_len"], 10, 64)
		if e != nil {
			return f, errors.New("invalid value for `max_len` option")
		}
		f.MaxLen(v, t.Options["max_len_err"])
	}
	if t.HasOption("pattern") {
		if !t.HasOption("pattern_err") {
			return f, errors.New("`pattern` option is set but `pattern_err` option is missing")
		}
		f.Pattern(t.Options["pattern"], t.Options["pattern_err"])
	}
	return f, nil
}

func parseTimeFieldOpts(f *DateTimeSchema, t *tagparser.Tag) (*DateTimeSchema, error) {
	if t.HasOption("epoch_unit") {
		switch t.Options["epoch_unit"] {
		case "nanosecond":
		case "ns":
			f.EpochUnit(EpochNanosecond)
		case "microsecond":
		case "us":
			f.EpochUnit(EpochMicrosecond)
		case "millisecond":
		case "ms":
			f.EpochUnit(EpochMillisecond)
		case "second":
		case "s":
			f.EpochUnit(EpochSecond)
		default:
			return f, errors.New("invalid value for `epoch_unit` option")
		}
	}
	if t.HasOption("layout") {
		f.Layout(t.Options["layout"])
	}
	return f, nil
}
