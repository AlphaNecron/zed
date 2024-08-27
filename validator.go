package zed

type (
	Validator struct {
		fields map[string]fieldTrait
	}
	OutValue struct {
		Field string
		Value any
	}
)

func New() *Validator {
	return &Validator{
		fields: make(map[string]fieldTrait),
	}
}

func (v *Validator) Field(name string, f fieldTrait) *Validator {
	v.fields[name] = f
	return v
}

func (v *Validator) Validate(m, out map[string]any) (e error) {
	for name, field := range v.fields {
		o, outOk := out[name]
		if !outOk {
			e = ErrOutFieldMissing
			return
		}
		if e = field.validate(m[name], o); e != nil {
			return e
		}
	}
	return
}
