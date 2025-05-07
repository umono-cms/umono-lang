package mocks

import (
	"regexp"

	"github.com/umono-cms/umono-lang/interfaces"
)

type Converter struct {
	converter func(string) string
}

func NewMockConverter(converter func(string) string) *Converter {
	return &Converter{
		converter: converter,
	}
}

func (m *Converter) Convert(umonoLang string) string {
	return m.converter(umonoLang)
}

func (*Converter) ConvertBuiltInComp(interfaces.Call) string {
	return ""
}

func BasicConverter(umonoLang string) string {
	re := regexp.MustCompile(`\s*input\s*`)
	return re.ReplaceAllString(umonoLang, "output")
}
