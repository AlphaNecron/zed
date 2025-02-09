package zed

import (
	"errors"
	"fmt"
	"github.com/ogen-go/ogen"
	"github.com/vmihailenco/tagparser/v2"
	"go.uber.org/multierr"
	"reflect"
	"strings"
)

var (
	_ Schema[any] = (*StructSchema[any])(nil)
	_ error       = (*StructCompileError)(nil)
)

type (
	StructSchema[T any] struct {
		err    error
		fields map[string]*structField
	}
	StructCompileError struct {
		Path string
		Err  error
	}
	structField struct {
		schema   schemaTrait
		required bool
	}
)

func newStructSchema(err string) *StructSchema[any] {
	return &StructSchema[any]{
		fields: make(map[string]*structField),
		err:    errors.New(err),
	}
}

func newStructSchemaFromType[T any](t reflect.Type, err string) (sf *StructSchema[T], e error) {
	sf = &StructSchema[T]{
		fields: make(map[string]*structField),
		err:    errors.New(err),
	}
	for i := range t.NumField() {
		f := t.Field(i)
		tag := tagparser.Parse(f.Tag.Get("zed"))
		var schema schemaTrait = nil
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
		case reflect.Struct:
			if isUuidType(f.Type) {
				schema = UUID(errMsg)
			} else if isTimeType(f.Type) {
				schema, e = parseTimeFieldOpts(DateTime(errMsg), tag)
			} else {
				schema, e = newStructSchemaFromType[any](f.Type, errMsg)
			}
		default:
			continue
		}
		sf.fields[actualName] = &structField{
			schema:   schema,
			required: tag.HasOption("required"),
		}
		if e != nil {
			e = newCompileErr(e, actualName)
			return
		}
	}
	return
}

func (f *StructSchema[T]) AddField(name string, schema schemaTrait, required bool) {
	f.fields[name] = &structField{
		schema:   schema,
		required: required,
	}
}

func (f *StructSchema[T]) Validate(m any, abortEarly bool) (out T, e error) {
	refType := reflect.TypeOf(m)
	refVal := reflect.ValueOf(m)
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
		refVal = refVal.Elem()
	}
	if refType.Kind() != reflect.Struct {
		e = f.err
		return
	}
	mp := make(map[string]reflect.Value)
	for i := range refType.NumField() {
		_f := refType.Field(i)
		name := _f.Name
		t := _f.Tag.Get("zed")
		tagName := t[:strings.Index(t, ",")]
		if tagName == "-" {
			continue
		}
		if tagName != "" {
			name = tagName
		}
		mp[name] = refVal.Field(i)
	}
	for name, field := range f.fields {
		refField, found := mp[name]
		if !found && !field.required {
			continue
		}
		var val any = nil
		if found {
			val = refField.Interface()
		}
		_, _e := field.schema.validateGeneric(val, abortEarly)
		if _e != nil {
			if abortEarly {
				e = newValidationErr(_e, name)
				return
			}
			e = multierr.Append(e, newValidationErr(_e, name))
		}
		// refField.Set(reflect.ValueOf(_out))
	}
	return
}

func (f *StructSchema[T]) validateGeneric(v any, abortEarly bool) (any, error) {
	return f.Validate(v, abortEarly)
}

func (f *StructSchema[T]) ToSchema() (s *ogen.Schema) {
	s = ogen.NewSchema().SetType("object")
	for name, field := range f.fields {
		prop := ogen.NewProperty().
			SetName(name).
			SetSchema(field.schema.ToSchema())
		if field.required {
			s.AddRequiredProperties(prop)
		} else {
			s.AddOptionalProperties(prop)
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
