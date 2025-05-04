package umonolang

type Parameter struct {
	name  string
	typ   string
	value any
}

func NewParam(name, typ string, value any) *Parameter {
	return &Parameter{
		name:  name,
		typ:   typ,
		value: value,
	}
}

func (p *Parameter) Name() string {
	return p.name
}

func (p *Parameter) Type() string {
	return p.typ
}

func (p *Parameter) Value() any {
	return p.value
}

func (p *Parameter) SetValue(val any) {
	p.value = val
}

func (p *Parameter) ValueAsString() string {
	if p.typ == "string" {
		return p.value.(string)
	}
	if p.typ == "bool" {
		if val := p.value.(bool); val {
			return "true"
		} else {
			return "false"
		}
	}
	return ""
}
