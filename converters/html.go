package converters

import (
	"bytes"

	"github.com/umono-cms/umono-lang/interfaces"
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

func (h *HTML) ConvertBuiltInComp(call interfaces.Call) string {
	if call.Component().Name() == "LINK" {
		return h.renderLink(call)
	} else if call.Component().Name() == "404" {
		return h.render404()
	}
	return ""
}

func (*HTML) renderLink(call interfaces.Call) string {
	text := call.ParameterByName("text")
	url := call.ParameterByName("url")
	newTab := call.ParameterByName("new-tab")

	if text == nil || url == nil || newTab == nil {
		// NOTE: Unexpected
		return ""
	}

	newTabStr := ""
	if newTab.Value().(bool) == true {
		newTabStr = ` target="_blank" rel="noopener noreferrer"`
	}

	return "<a href=\"" + url.Value().(string) + "\"" + newTabStr + ">" + text.Value().(string) + "</a>"
}

func (*HTML) render404() string {
	return "<h1>404 Not Found</h1>Page not found"
}
