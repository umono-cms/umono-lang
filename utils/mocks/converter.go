package mocks

import (
	"regexp"

	"github.com/umono-cms/umono-lang/interfaces"
)

type Converter struct {
}

func (m *Converter) Convert(umonoLang string) string {
	re := regexp.MustCompile(`\s*input\s*`)
	return re.ReplaceAllString(umonoLang, "output")
}

func (*Converter) ConvertBuiltInComp(interfaces.Call) string {
	return ""
}
