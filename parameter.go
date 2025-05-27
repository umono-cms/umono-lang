package umonolang

import (
	"strings"

	"github.com/umono-cms/umono-lang/interfaces"
	"github.com/umono-cms/umono-lang/internal/parameters"
	ustrings "github.com/umono-cms/umono-lang/internal/utils/strings"
)

func ParameterByName(params []interfaces.Parameter, name string) interfaces.Parameter {
	for _, param := range params {
		if param.Name() == "name" {
			return param
		}
	}
	return nil
}

func readParams(raw string, indexes [][]int) []interfaces.Parameter {

	if raw == "" {
		return []interfaces.Parameter{}
	}

	params := []interfaces.Parameter{}
	for i := 0; i < len(indexes); i++ {

		start := indexes[i][0]

		var end int
		if i != len(indexes)-1 {
			end = indexes[i+1][0]
		} else {
			newLineIndex := strings.Index(raw[start:], "\n")
			if newLineIndex != -1 {
				end = start + newLineIndex
			} else {
				end = ustrings.LastRuneIndex(raw) + 1
			}
		}

		ok, key, val := ustrings.SeparateKeyValue(raw[start:end], `\s*=\s*`, ``)
		if !ok || !validParamKey(key) {
			break
		}

		invalidValue, valAny, typ := getValueWithType(val)
		if !invalidValue {
			break
		}

		params = append(params, parameters.NewDynamicParam(key, typ, valAny))
	}

	return params
}

func validParamKey(key string) bool {

	valTrimmed := strings.TrimSpace(key)

	if valTrimmed[0] == '-' || valTrimmed[len(valTrimmed)-1] == '-' {
		return false
	}

	if strings.ContainsRune(valTrimmed, '_') {
		return false
	}

	valRune := []rune(valTrimmed)

	for i := 0; i < len(valRune)-1; i++ {
		if valRune[i] == '-' && valRune[i+1] == '-' {
			return false
		}
	}

	return true
}

func getValueWithType(val string) (bool, any, string) {

	val = strings.TrimSpace(val)

	if val == "" {
		return false, nil, ""
	}

	quoteIndexes := ustrings.FindAllStringIndex(val, `".*"`)

	if len(quoteIndexes) > 0 {
		return true, string([]rune(val)[1 : quoteIndexes[0][1]-1]), "string"
	}

	htmlQuoteIndexes := ustrings.FindAllStringIndex(val, "&quot;.*&quot;")

	if len(htmlQuoteIndexes) > 0 {
		return true, string([]rune(val)[6 : htmlQuoteIndexes[0][1]-6]), "string"
	}

	boolIndexes := ustrings.FindAllStringIndex(val, `true|false`)

	if len(boolIndexes) > 0 {

		val := string([]rune(val)[boolIndexes[0][0]:boolIndexes[0][1]])
		boolVal := false

		if val == "true" {
			boolVal = true
		}

		return true, boolVal, "bool"
	}

	return false, nil, ""
}
