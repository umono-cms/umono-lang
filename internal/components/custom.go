package components

import "github.com/umono-cms/umono-lang/interfaces"

type Custom struct {
	name       string
	rawContent string
	params     []interfaces.Parameter
}

func NewCustom(name, rawContent string) interfaces.Component {
	return &Custom{
		name:       name,
		rawContent: rawContent,
	}
}

func NewCustomWithParams(name, rawContent string, params []interfaces.Parameter) interfaces.Component {
	return &Custom{
		name:       name,
		rawContent: rawContent,
		params:     params,
	}
}

func (c *Custom) Name() string {
	return c.name
}

func (c *Custom) Parameters() []interfaces.Parameter {
	return c.params
}

func (c *Custom) RawContent() string {
	return c.rawContent
}

func (c *Custom) NeedToConvert() bool {
	return false
}
