package components

import (
	"github.com/umono-cms/umono-lang/interfaces"
	args "github.com/umono-cms/umono-lang/internal/arguments/link"
)

type Link struct{}

func (*Link) Name() string {
	return "LINK"
}

func (*Link) Arguments() []interfaces.Argument {
	return []interfaces.Argument{
		&args.URL{},
		&args.Text{},
		&args.NewTab{},
	}
}

func (*Link) RawContent() string  { return "" }
func (*Link) NeedToConvert() bool { return true }
