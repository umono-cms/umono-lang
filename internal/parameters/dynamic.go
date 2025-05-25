package parameters

import "github.com/umono-cms/umono-lang/interfaces"

type Dynamic struct {
	name  string
	typ   string
	deflt any
}

func NewDynamicParam(name, typ string, deflt any) interfaces.Parameter {
	return &Dynamic{
		name:  name,
		typ:   typ,
		deflt: deflt,
	}
}

func (d *Dynamic) Name() string {
	return d.name
}

func (d *Dynamic) Type() string {
	return d.typ
}

func (d *Dynamic) Default() any {
	return d.deflt
}
