package strings

import (
	"regexp"
	"strings"
)

func Indexes(s, substr string) []int {

	indexes := []int{}
	start := 0

	for {
		index := strings.Index(s[start:], substr)
		if index == -1 {
			break
		}
		realIndex := start + index
		indexes = append(indexes, realIndex)

		start = realIndex + 1
	}

	return indexes
}

func IndexesByRegex(s, regex string) []int {

	re := regexp.MustCompile(regex)
	matches := re.FindAllStringIndex(s, -1)

	indexes := []int{}
	for _, match := range matches {
		indexes = append(indexes, match[0])
	}

	return indexes
}

func FindAllString(s string, regex string, trimRegex string) []string {
	re := regexp.MustCompile(regex)
	strs := re.FindAllString(s, -1)

	if trimRegex == "" {
		return strs
	}

	re2 := regexp.MustCompile(trimRegex)

	trimmed := []string{}
	for _, str := range strs {
		trimmed = append(trimmed, re2.ReplaceAllString(str, ""))
	}

	return trimmed
}

func FindAllStringIndex(s string, regex string) [][]int {
	re := regexp.MustCompile(regex)
	indexes := re.FindAllStringIndex(s, -1)

	byteToRune := func(byteIdx int) int {
		return len([]rune(s[:byteIdx]))
	}

	runeIndexes := [][]int{}
	for _, idx := range indexes {
		runeIndexes = append(runeIndexes, []int{byteToRune(idx[0]), byteToRune(idx[1])})
	}

	return runeIndexes
}

func ReplaceSubstring(s string, newSub string, start int, end int) string {
	runes := []rune(s)

	if start < 0 || end > len(runes) || start > end {
		return s
	}

	return string(runes[:start]) + newSub + string(runes[end:])
}

func IsNumericScreamingSnakeCase(s string) bool {
	re := regexp.MustCompile(`^[A-Z0-9]+(?:_[A-Z0-9]+)*$`)
	return re.MatchString(s)
}
