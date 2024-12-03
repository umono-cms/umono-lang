package converters

import (
	"bytes"

	"github.com/yuin/goldmark"
)

type HTML struct {
}

func NewHTML() *HTML {
	return &HTML{}
}

func (*HTML) Convert(umonoLang string) string {

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(umonoLang), &buf); err != nil {
		return ""
	}

	html := buf.String()

	return html
}
