package mocks

import (
	"regexp"
)

type Converter struct {
}

func (m *Converter) Convert(umonoLang string) string {
	re := regexp.MustCompile(`\s*input\s*`)
	return re.ReplaceAllString(umonoLang, "output")
}
