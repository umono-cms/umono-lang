package converters

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

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

func (h *HTML) renderLink(call interfaces.Call) string {
	text := call.ArgumentByName("text")
	url := call.ArgumentByName("url")
	newTab := call.ArgumentByName("new-tab")

	if text == nil || url == nil || newTab == nil {
		// NOTE: Unexpected
		return ""
	}

	newTabStr := ""
	if newTab.Value().(bool) == true {
		newTabStr = ` target="_blank" rel="noopener noreferrer"`
	}

	return "<a href=\"" + h.convertAsInline(url.Value().(string)) + "\"" + newTabStr + ">" + h.convertAsInline(text.Value().(string)) + "</a>"
}

func (*HTML) render404() string {
	return "<h1>404 Not Found</h1>Page not found"
}

func (h *HTML) convertAsInline(content string) string {

	blocker := "blocker-to-prevent-block-elements"

	raw := blocker + content + blocker

	html := h.Convert(raw)

	pattern := fmt.Sprintf(`(?i)^<p>\s*%s\s*|\s*%s\s*</p>\s*$`, blocker, blocker)

	re := regexp.MustCompile(pattern)

	output := re.ReplaceAllString(html, "")

	if strings.Contains(output, "<!-- raw HTML omitted -->") {
		return content
	}

	return output
}
