package components

import (
	args "github.com/umono-cms/umono-lang/arguments/link"
	"github.com/umono-cms/umono-lang/interfaces"
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

func (*Link) RawContent() string    { return "" }
func (*Link) PutAfterConvert() bool { return true }
