package components

import "github.com/umono-cms/umono-lang/interfaces"

type S404 struct{}

func (*S404) Name() string {
	return "404"
}

func (*S404) Parameters() []interfaces.Parameter {
	return []interfaces.Parameter{}
}

func (*S404) RawContent() string  { return "" }
func (*S404) NeedToConvert() bool { return true }
