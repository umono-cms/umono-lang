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
	converter        interfaces.Converter
	globalCompMap    map[string]string
	builtInComps     []interfaces.Component
	builtInCompNames []string
}

func New(converter interfaces.Converter) *UmonoLang {

	builtInComps := []interfaces.Component{
		&components.Container{},
		&components.Row{},
		&components.Col{},
		&components.Link{},
	}

	builtInCompNames := []string{}
	for _, b := range builtInComps {
		builtInCompNames = append(builtInCompNames, b.Name())
	}

	return &UmonoLang{
		converter:        converter,
		globalCompMap:    make(map[string]string),
		builtInComps:     builtInComps,
		builtInCompNames: builtInCompNames,
	}
}

func (ul *UmonoLang) Convert(raw string) string {

	realContent := raw
	localeCompMap := map[string]string{}

	firstCompDefIndex := ul.findFirstCompDefIndex(raw)

	if firstCompDefIndex != -1 {
		realContent = raw[:firstCompDefIndex]
		localeCompMap = ul.readLocaleComponents(raw[firstCompDefIndex:])
	}

	// TODO: Complete it
	builtInCompMap := ul.readBuiltInComponents(realContent)

	compMap := builtInCompMap

	for name, content := range ul.globalCompMap {
		compMap[name] = content
	}

	for name, content := range localeCompMap {
		compMap[name] = content
	}

	return ul.handleComps(realContent, compMap, 1)
}

func (ul *UmonoLang) SetGlobalComponent(name, content string) error {

	if !ustrings.IsNumericScreamingSnakeCase(name) {
		return errors.New("SYNTAX_ERROR: Component names have to be SCREAMING_SNAKE_CASE.")
	}

	ul.globalCompMap[name] = content

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

	re := regexp.MustCompile(`\n~\s+[A-Z0-9]+(?:_[A-Z0-9]+)?\s*\n`)

	match := re.FindStringIndex(raw)

	if match != nil {
		return match[0]
	}

	return -1
}

func (ul *UmonoLang) readLocaleComponents(localeCompsRaw string) map[string]string {

	localeCompsIndexes := ustrings.IndexesByRegex(localeCompsRaw, `\n~\s+[A-Z0-9_]+(?:_[A-Z0-9]+)?\s*\n`)

	compContentMap := map[string]string{}

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

		compContentMap[re.ReplaceAllString(compNameRaw, "")] = strings.TrimSpace(compContentRaw)
	}

	return compContentMap
}

func (ul *UmonoLang) readBuiltInComponents(raw string) map[string]string {
	return map[string]string{}
}

func (ul *UmonoLang) handleComps(content string, compMap map[string]string, deep int) string {

	if deep == 20 {
		return ""
	}

	comps := ustrings.FindAllString(content, `\s*[A-Z0-9_]+(?:_[A-Z0-9]+)?\s*`, `^\s*|\s*$`)

	converted := ul.converter.Convert(content)

	for _, comp := range comps {
		cont, ok := compMap[comp]
		if !ok {
			continue
		}

		re := regexp.MustCompile(fmt.Sprintf(`%s`, comp))
		converted = re.ReplaceAllString(converted, ul.handleComps(cont, compMap, deep+1))
	}

	return strings.TrimSpace(converted)
}
