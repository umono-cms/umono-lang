package umonolang

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/umono-cms/umono-lang/components"
	"github.com/umono-cms/umono-lang/interfaces"
	ustrings "github.com/umono-cms/umono-lang/utils/strings"
)

type UmonoLang struct {
	converter   interfaces.Converter
	globalComps []interfaces.Component
}

func New(converter interfaces.Converter) *UmonoLang {
	return &UmonoLang{
		converter:   converter,
		globalComps: []interfaces.Component{},
	}
}

func (ul *UmonoLang) Convert(raw string) string {

	content := raw
	localComps := []interfaces.Component{}

	firstCompDefIndex := ul.findFirstCompDefIndex(raw)

	if firstCompDefIndex != -1 {
		content = raw[:firstCompDefIndex]
		localComps = readLocalComps(raw[firstCompDefIndex:])
	}

	comps := builtInComps()
	comps = overrideComps(comps, ul.globalComps)
	comps = overrideComps(comps, localComps)

	cursor := 0
	preConverted := ul.converter.Convert(ul.handleComps(comps, content, 1, cursor))

	return ul.convert(comps, preConverted)
}

func (ul *UmonoLang) SetGlobalComponent(name, content string) error {

	if !ustrings.IsNumericScreamingSnakeCase(name) {
		return errors.New("SYNTAX_ERROR: Component names have to be SCREAMING_SNAKE_CASE.")
	}

	ul.globalComps = overrideComps(ul.globalComps, []interfaces.Component{components.NewCustom(name, content)})

	return nil
}

func (ul *UmonoLang) RemoveGlobalComponent(name string) error {

	if !ustrings.IsNumericScreamingSnakeCase(name) {
		return errors.New("SYNTAX_ERROR: Component names have to be SCREAMING_SNAKE_CASE.")
	}

	index, found := findCompByName(ul.globalComps, name)
	if found == nil {
		return fmt.Errorf("NOT_FOUND: The global component named '%s' not found.", name)
	}

	ul.globalComps = append(ul.globalComps[:index], ul.globalComps[index+1:]...)

	return nil
}

func (ul *UmonoLang) findFirstCompDefIndex(raw string) int {

	re := regexp.MustCompile(`\n~\s+[A-Z0-9]+(?:_[A-Z0-9]+)*\s*\n`)

	match := re.FindStringIndex(raw)

	if match != nil {
		return match[0]
	}

	return -1
}

func (ul *UmonoLang) handleComps(comps []interfaces.Component, content string, deep int, cursor int) string {

	if deep == 20 {
		return ""
	}

	calls := readCalls(content, comps)

	handled := content

	for _, call := range calls {
		if call.Component().NeedToConvert() {
			continue
		}

		handled = ustrings.ReplaceSubstring(handled, ul.handleComps(comps, call.Component().RawContent(), deep+1, cursor), call.Start()+cursor, call.End()+cursor)

		rawContentLen := utf8.RuneCountInString(call.Component().RawContent())
		callLen := call.End() - call.Start()

		abs := rawContentLen - callLen
		if abs < 0 {
			abs = -abs
		}

		if callLen < rawContentLen {
			cursor += abs
		} else {
			cursor -= abs
		}

	}

	return strings.TrimSpace(handled)
}

func (ul *UmonoLang) convert(comps []interfaces.Component, handled string) string {

	converted := handled

	calls := readCalls(converted, comps)

	cursor := 0

	for _, call := range calls {
		if !call.Component().NeedToConvert() {
			continue
		}

		output := ul.converter.ConvertBuiltInComp(call)

		converted = ustrings.ReplaceSubstring(converted, output, call.Start()+cursor, call.End()+cursor)

		convertedLen := utf8.RuneCountInString(output)
		callLen := call.End() - call.Start()

		abs := convertedLen - callLen
		if abs < 0 {
			abs = -abs
		}

		if callLen < convertedLen {
			cursor += abs
		} else {
			cursor -= abs
		}
	}

	return converted
}
