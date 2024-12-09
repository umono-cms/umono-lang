package umonolang

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/umono-cms/umono-lang/interfaces"
	ustrings "github.com/umono-cms/umono-lang/utils/strings"
)

type UmonoLang struct {
	converter interfaces.Converter
}

func New(converter interfaces.Converter) *UmonoLang {
	return &UmonoLang{
		converter: converter,
	}
}

func (ul *UmonoLang) Convert(raw string) string {

	realContent := raw
	localeCompMaps := map[string]string{}

	firstCompDefIndex := ul.findFirstCompDefIndex(raw)

	if firstCompDefIndex != -1 {
		realContent = raw[:firstCompDefIndex]
		localeCompMaps = ul.readLocaleComponents(raw[firstCompDefIndex:])
	}

	realContent = ul.convert(realContent, localeCompMaps, 1)

	return realContent
}

func (ul *UmonoLang) findFirstCompDefIndex(raw string) int {

	re := regexp.MustCompile(`\n~\s*[A-Z0-9_]+\s*\n`)

	match := re.FindStringIndex(raw)

	if match != nil {
		return match[0]
	}

	return -1
}

func (ul *UmonoLang) readLocaleComponents(localeCompsRaw string) map[string]string {

	localeCompsIndexes := ustrings.IndexesByRegex(localeCompsRaw, `\n~\s*[A-Z0-9_]+\s*\n`)

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

func (ul *UmonoLang) convert(content string, compMap map[string]string, deep int) string {

	if deep == 20 {
		return ""
	}

	comps := ustrings.FindAllString(content, `\{\{\s*[A-Z0-9_]+\s*\}\}`, `^\s*\{\{\s*|\s*\}\}\s*$`)

	contConverted := ul.converter.Convert(content)

	for _, comp := range comps {
		cont, ok := compMap[comp]
		if !ok {
			continue
		}

		converted := ul.convert(cont, compMap, deep+1)
		re := regexp.MustCompile(fmt.Sprintf(`\{\{\s*%s\s*\}\}`, comp))
		contConverted = re.ReplaceAllString(contConverted, "\n"+converted+"\n")
	}

	return contConverted
}
