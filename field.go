package zed

import "math"

const (
	fieldString  fieldKind = "string"
	fieldBool    fieldKind = "bool"
	fieldUuid    fieldKind = "uuid"
	fieldUint8   fieldKind = "uint8"
	fieldUint16  fieldKind = "uint16"
	fieldUint32  fieldKind = "uint32"
	fieldInt8    fieldKind = "int8"
	fieldInt16   fieldKind = "int16"
	fieldInt32   fieldKind = "int32"
	fieldInt64   fieldKind = "int64"
	fieldFloat32 fieldKind = "float32"
	fieldFloat64 fieldKind = "float64"
)

type (
	fieldKind string
)

func String() *StringField {
	return newStrField()
}

func Uint8(err string) *NumField[uint8] {
	return newNumField[uint8](0, math.MaxUint8, err)
}

func Uint16(err string) *NumField[uint16] {
	return newNumField[uint16](0, math.MaxUint16, err)
}

func Uint32(err string) *NumField[uint32] {
	return newNumField[uint32](0, math.MaxUint32, err)
}

func Int8(err string) *NumField[int8] {
	return newNumField[int8](math.MinInt8, math.MaxInt8, err)
}

func Int16(err string) *NumField[int16] {
	return newNumField[int16](math.MinInt16, math.MaxInt16, err)
}

func Int32(err string) *NumField[int32] {
	return newNumField[int32](math.MinInt32, math.MaxInt32, err)
}

func Int64(err string) *NumField[int64] {
	return newNumField[int64](math.MinInt64, math.MaxInt64, err)
}

func Float32(err string) *NumField[float32] {
	return newNumField[float32](0, 0, err)
}

func Float64(err string) *NumField[float64] {
	return newNumField[float64](0, 0, err)
}

func Int(err string) *NumField[int] {
	return newNumField[int](math.MinInt, math.MaxInt, err)
}
