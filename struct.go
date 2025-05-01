package zed

import (
	"errors"
	"fmt"
	"github.com/ogen-go/ogen"
	"github.com/vmihailenco/tagparser/v2"
	"go.uber.org/multierr"
	"reflect"
)

var (
	_ Schema[any] = (*StructSchema[any])(nil)
	_ error       = (*StructCompileError)(nil)
)

type (
	StructSchema[T any] struct {
		*baseSchema[T, T]
		fields map[string]*structField
	}
	StructCompileError struct {
		Path string
		Err  error
	}
	structField struct {
		schema SchemaTrait
		flags  StructFieldFlag
	}
	StructFieldFlag uint8
)

const (
	FieldRequired StructFieldFlag = 1 << iota
)

func (f StructFieldFlag) Has(_f StructFieldFlag) bool {
	return f&_f != 0
}

func newStructSchema[T any](err string) *StructSchema[T] {
	return &StructSchema[T]{
		fields:     make(map[string]*structField),
		baseSchema: newBaseSchema[T, T](err),
	}
}

func newStructSchemaFromType[T any](t reflect.Type, err string) (sf *StructSchema[T], e error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct && t.Kind() != reflect.Map {
		e = errors.New("expected out to be `map` or `struct`")
	}
	sf = newStructSchema[T](err)
	for i := range t.NumField() {
		f := t.Field(i)
		tag := tagparser.Parse(f.Tag.Get("zed"))
		var schema SchemaTrait = nil
		actualName := f.Name
		if tag.Name == "-" {
			continue
		}
		if tag.Name != "" {
			actualName = tag.Name
		}
		if !tag.HasOption("err") {
			e = newCompileErr(errors.New("missing `err` option"), actualName)
			return
		}
		errMsg := tag.Options["err"]
		switch f.Type.Kind() {
		case reflect.String:
			schema, e = parseStringFieldOpts(String(errMsg), tag)
		case reflect.Bool:
			bs := Bool(errMsg)
			// if tag.HasOption("strict") {
			// 	bf.Strict()
			// }
			schema = bs
		case reflect.Int:
			schema, e = parseNumFieldOpts(Int(errMsg), tag)
		case reflect.Int8:
			schema, e = parseNumFieldOpts(Int8(errMsg), tag)
		case reflect.Int16:
			schema, e = parseNumFieldOpts(Int16(errMsg), tag)
		case reflect.Int32:
			schema, e = parseNumFieldOpts(Int32(errMsg), tag)
		case reflect.Int64:
			schema, e = parseNumFieldOpts(Int64(errMsg), tag)
		case reflect.Uint8:
			schema, e = parseNumFieldOpts(Uint8(errMsg), tag)
		case reflect.Uint16:
			schema, e = parseNumFieldOpts(Uint16(errMsg), tag)
		case reflect.Uint32:
			schema, e = parseNumFieldOpts(Uint32(errMsg), tag)
		case reflect.Float32:
			schema, e = parseNumFieldOpts(Float32(errMsg), tag)
		case reflect.Float64:
			schema, e = parseNumFieldOpts(Float64(errMsg), tag)
		case reflect.Array:
			if isUuidType(f.Type) {
				schema = UUID(errMsg)
			}
		case reflect.Struct:
			if isTimeType(f.Type) {
				schema, e = parseTimeFieldOpts(DateTime(errMsg), tag)
			} else {
				schema, e = newStructSchemaFromType[any](f.Type, errMsg)
			}
		default:
			continue
		}
		var flags StructFieldFlag = 0
		if tag.HasOption("required") {
			flags |= FieldRequired
		}
		sf.fields[actualName] = &structField{
			schema: schema,
			flags:  flags,
		}
		if e != nil {
			e = newCompileErr(e, actualName)
			return
		}
	}
	return
}

func (s *StructSchema[T]) AddField(name string, schema SchemaTrait, flags StructFieldFlag) *StructSchema[T] {
	s.fields[name] = &structField{
		schema: schema,
		flags:  flags,
	}
	return s
}

func (s *StructSchema[T]) Validate(m any, flags SchemaValidationFlag) (out T, e error) {
	var refOut reflect.Value
	refVal := reflect.ValueOf(m)
	refType := refVal.Type()
	outType := reflect.TypeFor[T]()
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
		refVal = refVal.Elem()
	}
	if outType.Kind() == reflect.Ptr {
		outType = outType.Elem()
	}
	if outType.Kind() == reflect.Struct {
		refOut = reflect.New(outType).Elem()
	} else if outType.Kind() == reflect.Map {
		refOut = reflect.MakeMap(outType)
	}
	mp := make(map[string]reflect.Value)
	if refType.Kind() == reflect.Map {
		if refType.Key().Kind() != reflect.String {
			e = s.err
			return
		}
		for _, k := range refVal.MapKeys() {
			mp[k.String()] = refVal.MapIndex(k)
		}
	} else if refType.Kind() == reflect.Struct {
		for i := range refType.NumField() {
			if name, ok := parseStructFieldName(refType.Field(i)); ok {
				mp[name] = refVal.Field(i)
			}
		}
	} else {
		e = s.err
		return
	}
	structNameMp := make(map[string]string)
	if outType.Kind() == reflect.Struct {
		for i := range outType.NumField() {
			f := outType.Field(i)
			if name, ok := parseStructFieldName(f); ok {
				structNameMp[name] = f.Name
			}
		}
	}
	for name, field := range s.fields {
		refField, found := mp[name]
		if !found && !field.flags.Has(FieldRequired) {
			continue
		}
		var val any = nil
		if found {
			val = refField.Interface()
		}
		_out, _e := field.schema.ValidateGeneric(val, flags)
		if _e != nil {
			if flags.Has(AbortEarly) {
				e = newValidationErr(_e, name)
				return
			}
			e = multierr.Append(e, newValidationErr(_e, name))
			continue
		}
		if outType.Kind() == reflect.Map {
			refOut.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(_out))
		} else if sname, ok := structNameMp[name]; ok {
			f := refOut.FieldByName(sname)
			if f.CanSet() {
				f.Set(reflect.ValueOf(_out))
			}
		}
	}
	out = refOut.Interface().(T)
	return
}

func (s *StructSchema[T]) ValidateGeneric(v any, flags SchemaValidationFlag) (any, error) {
	return s.Validate(v, flags)
}

func (s *StructSchema[T]) ToSchema() (os *ogen.Schema) {
	os = ogen.NewSchema().SetType("object")
	for name, field := range s.fields {
		prop := ogen.NewProperty().
			SetName(name).
			SetSchema(field.schema.ToSchema())
		if field.flags.Has(FieldRequired) {
			os.AddRequiredProperties(prop)
		} else {
			os.AddOptionalProperties(prop)
		}
	}
	return
}

func (fe *StructCompileError) Error() string {
	return fmt.Sprintf("field `%s`: %v", fe.Path, fe.Err)
}

func newCompileErr(err error, field string) error {
	if err == nil {
		return nil
	}
	var fe *StructCompileError
	if errors.As(err, &fe) {
		if fe.Path == "" {
			fe.Path = field
		} else {
			fe.Path = fe.Path + "." + field
		}
	} else {
		fe = &StructCompileError{
			Path: field,
			Err:  err,
		}
	}
	return fe
}
