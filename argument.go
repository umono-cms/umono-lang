package umonolang

import (
	"strings"

	"github.com/umono-cms/umono-lang/arguments"
	"github.com/umono-cms/umono-lang/interfaces"
	ustrings "github.com/umono-cms/umono-lang/utils/strings"
)

func readArgs(raw string, indexes [][]int) []interfaces.Argument {

	if raw == "" {
		return []interfaces.Argument{}
	}

	args := []interfaces.Argument{}
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
		if !ok || !validArgKey(key) {
			break
		}

		invalidValue, valAny, typ := getValueWithType(val)
		if !invalidValue {
			break
		}

		args = append(args, arguments.NewDynamicArg(key, typ, valAny))
	}

	return args
}

func validArgKey(key string) bool {

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

	if val == "true" || val == "false" {
		return true, val, "bool"
	}

	if (val[0] == '"' && val[ustrings.LastRuneIndex(val)] == '"') || (val[0] == '\'' && val[ustrings.LastRuneIndex(val)] == '\'') {

		return true, val[1:ustrings.LastRuneIndex(val)], "string"
	}

	return false, nil, ""
}
