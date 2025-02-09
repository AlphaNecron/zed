package zed

import (
	"math"
	"reflect"
)

type Schema[TOut any] interface {
	Validate(val any, abortEarly bool) (TOut, error)
	schemaTrait
}

func Bool(err string) *BoolSchema {
	return newBoolSchema(err)
}

func String(err string) *StringSchema {
	return newStrSchema(err)
}

func Uint8(err string) *NumSchema[uint8] {
	return newNumSchema[uint8](0, math.MaxUint8, err)
}

func Uint16(err string) *NumSchema[uint16] {
	return newNumSchema[uint16](0, math.MaxUint16, err)
}

func Uint32(err string) *NumSchema[uint32] {
	return newNumSchema[uint32](0, math.MaxUint32, err)
}

func Int8(err string) *NumSchema[int8] {
	return newNumSchema[int8](math.MinInt8, math.MaxInt8, err)
}

func Int16(err string) *NumSchema[int16] {
	return newNumSchema[int16](math.MinInt16, math.MaxInt16, err)
}

func Int32(err string) *NumSchema[int32] {
	return newNumSchema[int32](math.MinInt32, math.MaxInt32, err)
}

func Int64(err string) *NumSchema[int64] {
	return newNumSchema[int64](math.MinInt64, math.MaxInt64, err)
}

func Float32(err string) *NumSchema[float32] {
	return newNumSchema[float32](0, 0, err)
}

func Float64(err string) *NumSchema[float64] {
	return newNumSchema[float64](0, 0, err)
}

func Int(err string) *NumSchema[int] {
	return newNumSchema[int](math.MinInt, math.MaxInt, err)
}

func UUID(err string) *UUIDSchema {
	return newUuidSchema(err)
}

func DateTime(err string) *DateTimeSchema {
	return newDateTimeSchema(err)
}

func StructFor[T any](err string) *StructSchema[T] {
	f, e := StructForE[T](err)
	if e != nil {
		panic(e)
	}
	return f
}

func StructForE[T any](err string) (*StructSchema[T], error) {
	return newStructSchemaFromType[T](reflect.TypeFor[T](), err)
}

func Struct(err string) *StructSchema[any] {
	return newStructSchema(err)
}
