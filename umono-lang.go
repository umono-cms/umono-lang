package umonolang

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/umono-cms/umono-lang/components"
	"github.com/umono-cms/umono-lang/interfaces"
	ustrings "github.com/umono-cms/umono-lang/utils/strings"
)

type UmonoLang struct {
	converter     interfaces.Converter
	globalCompMap map[string]interfaces.Component
}

func New(converter interfaces.Converter) *UmonoLang {
	return &UmonoLang{
		converter:     converter,
		globalCompMap: make(map[string]interfaces.Component),
	}
}

func (ul *UmonoLang) Convert(raw string) string {

	content := raw
	localCompMap := []interfaces.Component{}

	firstCompDefIndex := ul.findFirstCompDefIndex(raw)

	if firstCompDefIndex != -1 {
		content = raw[:firstCompDefIndex]
		localCompMap = ul.readLocalComponents(raw[firstCompDefIndex:])
	}

	comps := localCompMap

	/*compMap := builtInCompMap()

	for name, gc := range ul.globalCompMap {
		compMap[name] = gc
	}

	for name, lc := range localCompMap {
		compMap[name] = lc
	}*/

	cursor := 0
	return ul.converter.Convert(ul.handleComps(comps, content, 1, cursor))
}

func (ul *UmonoLang) SetGlobalComponent(name, content string) error {

	if !ustrings.IsNumericScreamingSnakeCase(name) {
		return errors.New("SYNTAX_ERROR: Component names have to be SCREAMING_SNAKE_CASE.")
	}

	ul.globalCompMap[name] = components.NewCustom(name, content)

	return nil
}

func (ul *UmonoLang) RemoveGlobalComponent(name string) error {

	if !ustrings.IsNumericScreamingSnakeCase(name) {
		return errors.New("SYNTAX_ERROR: Component names have to be SCREAMING_SNAKE_CASE.")
	}

	_, ok := ul.globalCompMap[name]
	if !ok {
		return fmt.Errorf("NOT_FOUND: The global component named '%s' not found.", name)
	}

	delete(ul.globalCompMap, name)

	return nil
}

func builtInCompMap() map[string]interfaces.Component {
	bcm := map[string]interfaces.Component{}

	builtInComps := []interfaces.Component{
		&components.Link{},
	}

	for _, bc := range builtInComps {
		bcm[bc.Name()] = bc
	}

	return bcm
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

		abs := len(call.Component().RawContent()) - len(call.Component().Name())
		if abs < 0 {
			abs = -abs
		}

		if len(call.Component().Name()) < len(call.Component().RawContent()) {
			cursor += abs
		} else {
			cursor -= abs
		}

	}

	return strings.TrimSpace(handled)
}
