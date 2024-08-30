package zed

import "math"

func Bool(err string, strict bool) *BoolField {
	return newBoolField(err, strict)
}

func String(err string) *StringField {
	return newStrField(err)
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

func UUID(err string) *UUIDField {
	return newUuidField(err)
}
