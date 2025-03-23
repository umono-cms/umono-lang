package components

import "github.com/umono-cms/umono-lang/interfaces"

type S404 struct{}

func (*S404) Name() string {
	return "404"
}

func (*S404) Arguments() []interfaces.Argument {
	return []interfaces.Argument{}
}

func (*S404) RawContent() string  { return "" }
func (*S404) NeedToConvert() bool { return true }
