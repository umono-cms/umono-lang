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
	compContentMap := map[string]string{}

	firstTildeIndex := strings.Index(raw, "\n~")
	if firstTildeIndex != -1 {
		realContent = raw[:firstTildeIndex]
		compContentMap = ul.resolveCompContentMap(raw[firstTildeIndex:])
	}

	convertedRealContent := ul.converter.Convert(realContent)

	for compName, content := range compContentMap {
		re := regexp.MustCompile(fmt.Sprintf(`\{\{\s*%s\s*\}\}`, compName))
		convertedRealContent = re.ReplaceAllString(convertedRealContent, "\n"+ul.converter.Convert(content)+"\n")
	}

	return convertedRealContent
}

func (ul *UmonoLang) resolveCompContentMap(localeCompsRaw string) map[string]string {

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

		compNameRaw := trimmed[0:endOfCompName]
		compContentRaw := trimmed[endOfCompName:]

		compContentMap[re.ReplaceAllString(compNameRaw, "")] = strings.TrimSpace(compContentRaw)
	}

	return compContentMap
}
