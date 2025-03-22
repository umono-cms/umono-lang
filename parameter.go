package umonolang

type Parameter struct {
	name  string
	value any
}

func NewParam(name string, value any) *Parameter {
	return &Parameter{
		name:  name,
		value: value,
	}
}

func (p *Parameter) Name() string {
	return p.name
}

func (p *Parameter) Value() any {
	return p.value
}

func (p *Parameter) SetValue(val any) {
	p.value = val
}
