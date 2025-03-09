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
		localComps = ul.readLocalComponents(raw[firstCompDefIndex:])
	}

	comps := builtInComps()
	comps = overrideComps(comps, ul.globalComps)
	comps = overrideComps(comps, localComps)

	cursor := 0
	return ul.converter.Convert(ul.handleComps(comps, content, 1, cursor))
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

func (ul *UmonoLang) readLocalComponents(localeCompsRaw string) []interfaces.Component {

	localeCompsIndexes := ustrings.IndexesByRegex(localeCompsRaw, `\n~\s+[A-Z0-9_]+(?:_[A-Z0-9]+)*\s*\n`)

	comps := []interfaces.Component{}

	re := regexp.MustCompile(`(?s)^~\s*|\s*\n$`)

	for i := 0; i < len(localeCompsIndexes); i++ {
		var compRaw string
		if i == len(localeCompsIndexes)-1 {
			compRaw = localeCompsRaw[localeCompsIndexes[i]:]
		} else {
			compRaw = localeCompsRaw[localeCompsIndexes[i]:localeCompsIndexes[i+1]]
		}

		trimmed := strings.TrimSpace(compRaw)
		endOfCompName := strings.Index(trimmed, "\n")

		var compNameRaw string
		var compContentRaw string

		if endOfCompName == -1 {
			compNameRaw = trimmed
			compContentRaw = ""
		} else {
			compNameRaw = trimmed[0:endOfCompName]
			compContentRaw = trimmed[endOfCompName:]
		}

		compName := re.ReplaceAllString(compNameRaw, "")
		comps = append(comps, components.NewCustom(compName, strings.TrimSpace(compContentRaw)))
	}

	return comps
}

func (ul *UmonoLang) readBuiltInComponents(raw string) map[string]string {
	return map[string]string{}
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
		nameLen := utf8.RuneCountInString(call.Component().Name())

		abs := rawContentLen - nameLen
		if abs < 0 {
			abs = -abs
		}

		if nameLen < rawContentLen {
			cursor += abs
		} else {
			cursor -= abs
		}

	}

	return strings.TrimSpace(handled)
}
