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
