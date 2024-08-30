package zed

type (
	Zed struct {
		fields map[string]fieldTrait
	}
	OutValue struct {
		Field string
		Value any
	}
)

func New() *Zed {
	return &Zed{
		fields: make(map[string]fieldTrait),
	}
}

func (z *Zed) Field(name string, f fieldTrait) *Zed {
	z.fields[name] = f
	return z
}

func (z *Zed) Validate(m, out map[string]any) (e error) {
	for name, field := range z.fields {
		o, outOk := out[name]
		if !outOk {
			e = ErrOutFieldMissing
			return
		}
		if e = field.validate(m[name], o); e != nil {
			return NewError(e, name)
		}
	}
	return
}
