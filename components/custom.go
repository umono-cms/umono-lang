package components

import "github.com/umono-cms/umono-lang/interfaces"

type Custom struct {
	name       string
	rawContent string
}

func NewCustom(name, rawContent string) interfaces.Component {
	return &Custom{
		name:       name,
		rawContent: rawContent,
	}
}

func (c *Custom) Name() string {
	return c.name
}

func (c *Custom) Arguments() []interfaces.Argument {
	return []interfaces.Argument{}
}

func (c *Custom) RawContent() string {
	return c.rawContent
}

func (c *Custom) PutAfterConvert() bool {
	return false
}
