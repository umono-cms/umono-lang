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

	realContent := raw
	localCompMap := map[string]interfaces.Component{}

	firstCompDefIndex := ul.findFirstCompDefIndex(raw)

	if firstCompDefIndex != -1 {
		realContent = raw[:firstCompDefIndex]
		localCompMap = ul.readLocalComponents(raw[firstCompDefIndex:])
	}

	compMap := builtInCompMap()

	for name, gc := range ul.globalCompMap {
		compMap[name] = gc
	}

	for name, lc := range localCompMap {
		compMap[name] = lc
	}

	preConverted := ul.converter.Convert(ul.handleComps(realContent, compMap, 1))

	return preConverted
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

func (ul *UmonoLang) findFirstCompDefIndex(raw string) int {

	re := regexp.MustCompile(`\n~\s+[A-Z0-9]+(?:_[A-Z0-9]+)*\s*\n`)

	match := re.FindStringIndex(raw)

	if match != nil {
		return match[0]
	}

	return -1
}

func (ul *UmonoLang) readLocalComponents(localeCompsRaw string) map[string]interfaces.Component {

	localeCompsIndexes := ustrings.IndexesByRegex(localeCompsRaw, `\n~\s+[A-Z0-9_]+(?:_[A-Z0-9]+)*\s*\n`)

	compContentMap := map[string]interfaces.Component{}

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
		compContentMap[compName] = components.NewCustom(compName, strings.TrimSpace(compContentRaw))
	}

	return compContentMap
}

func (ul *UmonoLang) readBuiltInComponents(raw string) map[string]string {
	return map[string]string{}
}

func (ul *UmonoLang) handleComps(content string, compMap map[string]interfaces.Component, deep int) string {

	if deep == 20 {
		return ""
	}

	comps := ustrings.FindAllString(content, `\s*[A-Z0-9_]+(?:_[A-Z0-9]+)*\s*`, `^\s*|\s*$`)

	handled := content

	for _, compName := range comps {
		comp, ok := compMap[compName]
		if !ok || comp.PutAfterConvert() {
			continue
		}

		re := regexp.MustCompile(fmt.Sprintf(`%s`, compName))
		handled = re.ReplaceAllString(handled, ul.handleComps(comp.RawContent(), compMap, deep+1))
	}

	return strings.TrimSpace(handled)
}
