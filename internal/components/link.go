package components

import (
	"github.com/umono-cms/umono-lang/interfaces"
	params "github.com/umono-cms/umono-lang/internal/parameters/link"
)

type Link struct{}

func (*Link) Name() string {
	return "LINK"
}

func (*Link) Parameters() []interfaces.Parameter {
	return []interfaces.Parameter{
		&params.URL{},
		&params.Text{},
		&params.NewTab{},
	}
}

func (*Link) RawContent() string  { return "" }
func (*Link) NeedToConvert() bool { return true }
