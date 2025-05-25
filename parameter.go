package umonolang

type parameter struct {
	name  string
	typ   string
	value any
}

func newParam(name, typ string, value any) *parameter {
	return &parameter{
		name:  name,
		typ:   typ,
		value: value,
	}
}

func (p *parameter) Name() string {
	return p.name
}

func (p *parameter) Type() string {
	return p.typ
}

func (p *parameter) Value() any {
	return p.value
}

func (p *parameter) SetValue(val any) {
	p.value = val
}

func (p *parameter) ValueAsString() string {
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
