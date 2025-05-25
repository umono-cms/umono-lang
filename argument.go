package umonolang

type argument struct {
	name  string
	typ   string
	value any
}

func newArg(name, typ string, value any) *argument {
	return &argument{
		name:  name,
		typ:   typ,
		value: value,
	}
}

func (a *argument) Name() string {
	return a.name
}

func (a *argument) Type() string {
	return a.typ
}

func (a *argument) Value() any {
	return a.value
}

func (a *argument) SetValue(val any) {
	a.value = val
}

func (a *argument) ValueAsString() string {
	if a.typ == "string" {
		return a.value.(string)
	}
	if a.typ == "bool" {
		if val := a.value.(bool); val {
			return "true"
		} else {
			return "false"
		}
	}
	return ""
}
