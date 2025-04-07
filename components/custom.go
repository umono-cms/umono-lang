package components

import "github.com/umono-cms/umono-lang/interfaces"

type Custom struct {
	name       string
	rawContent string
	args       []interfaces.Argument
}

func NewCustom(name, rawContent string) interfaces.Component {
	return &Custom{
		name:       name,
		rawContent: rawContent,
	}
}

func NewCustomWithArgs(name, rawContent string, args []interfaces.Argument) interfaces.Component {
	return &Custom{
		name:       name,
		rawContent: rawContent,
		args:       args,
	}
}

func (c *Custom) Name() string {
	return c.name
}

func (c *Custom) Arguments() []interfaces.Argument {
	return c.args
}

func (c *Custom) RawContent() string {
	return c.rawContent
}

func (c *Custom) NeedToConvert() bool {
	return false
}
