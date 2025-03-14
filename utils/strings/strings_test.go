package strings

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StringsTestSuite struct {
	suite.Suite
}

func (s *StringsTestSuite) TestIndexes() {
	for sI, scene := range []struct {
		str    string
		subStr string
		result []int
	}{
		{
			"Hello World",
			"Hello",
			[]int{0},
		},
		{
			"xoxoxo",
			"x",
			[]int{0, 2, 4},
		},
		{
			"Another Love",
			"love",
			[]int{},
		},
		{
			"Umono Lang",
			"g",
			[]int{9},
		},
	} {
		indexes := Indexes(scene.str, scene.subStr)
		require.Equal(s.T(), len(scene.result), len(indexes), "scene index: "+strconv.Itoa(sI))
		for i, _ := range scene.result {
			require.Equal(s.T(), scene.result[i], indexes[i], "scene index: "+strconv.Itoa(sI))
		}
	}
}

func (s *StringsTestSuite) TestIndexesByRegex() {
	for sI, scene := range []struct {
		str    string
		regex  string
		result []int
	}{
		{
			"Hello World",
			`\s`,
			[]int{5},
		},
		{
			"U m o n o",
			`\s`,
			[]int{1, 3, 5, 7},
		},
		{
			"Lang",
			`\s`,
			[]int{},
		},
		{
			"UmonoLang",
			`Umono`,
			[]int{0},
		},
	} {
		indexes := IndexesByRegex(scene.str, scene.regex)
		require.Equal(s.T(), len(scene.result), len(indexes), "scene index: "+strconv.Itoa(sI))
		for i, _ := range scene.result {
			require.Equal(s.T(), scene.result[i], indexes[i], "scene index: "+strconv.Itoa(sI))
		}
	}
}

func (s *StringsTestSuite) TestFindAllString() {
	for sI, scene := range []struct {
		str       string
		regex     string
		trimRegex string
		result    []string
	}{
		{
			"a HELLO_WORLD b",
			`\s*[A-Z0-9_]+\s*`,
			`^\s*|\s*$`,
			[]string{"HELLO_WORLD"},
		},
		{
			"a HELLO_WORLD b HELLO_ANOTHER_WORLD UMONO_LANG",
			`\s*[A-Z0-9_]+\s*`,
			`^\s*|\s*$`,
			[]string{"HELLO_WORLD", "HELLO_ANOTHER_WORLD", "UMONO_LANG"},
		},
		{
			"123umono-lang123u-m-o-n-o456",
			`[a-z-]+`,
			`[0-9]+`,
			[]string{"umono-lang", "u-m-o-n-o"},
		},
		{
			"a HELLO b",
			`HELLO`,
			``,
			[]string{"HELLO"},
		},
	} {
		strs := FindAllString(scene.str, scene.regex, scene.trimRegex)
		require.Equal(s.T(), len(scene.result), len(strs), "scene index: "+strconv.Itoa(sI))
		for i, _ := range scene.result {
			require.Equal(s.T(), scene.result[i], strs[i], "scene index: "+strconv.Itoa(sI))
		}
	}
}

func (s *StringsTestSuite) TestFindAllStringIndex() {
	for _, scene := range []struct {
		str    string
		regex  string
		result [][]int
	}{
		{
			"UMONO",
			"M",
			[][]int{[]int{1, 2}},
		},
		{
			"UMONO",
			"O",
			[][]int{
				[]int{2, 3},
				[]int{4, 5},
			},
		},
		{
			"UMONO",
			"T",
			[][]int{},
		},
		{
			"HELLO_WORLD",
			`WORLD`,
			[][]int{[]int{6, 11}},
		},
		// UTF-8
		{
			"UMONO öçğüşiı UMONO",
			`UMONO`,
			[][]int{
				[]int{0, 5},
				[]int{14, 19},
			},
		},
		{
			"UMONO öçğüşiı UMONO",
			`[öçğüşiı]+`,
			[][]int{[]int{6, 13}},
		},
	} {
		result := FindAllStringIndex(scene.str, scene.regex)
		require.Equal(s.T(), scene.result, result)
	}
}

func (s *StringsTestSuite) TestSeparateKeyValue() {
	for i, scene := range []struct {
		str   string
		sep   string
		ok    bool
		key   string
		value string
	}{
		{
			`A # B`,
			`\s*#\s*`,
			true,
			`A`,
			`B`,
		},
		{
			`A ö B`,
			`\s*ö\s*`,
			true,
			`A`,
			`B`,
		},
		{
			`öçöç #### şişi`,
			`\s*####\s*`,
			true,
			`öçöç`,
			`şişi`,
		},
		{
			`  A === B  `,
			`\s*===\s*`,
			true,
			`A`,
			`B`,
		},
	} {
		ok, key, value := SeparateKeyValue(scene.str, scene.sep)
		require.Equal(s.T(), scene.ok, ok, "index: "+strconv.Itoa(i))
		require.Equal(s.T(), scene.key, key, "index: "+strconv.Itoa(i))
		require.Equal(s.T(), scene.value, value, "index: "+strconv.Itoa(i))
	}
}

func (s *StringsTestSuite) TestReplaceSubstring() {
	for _, scene := range []struct {
		str    string
		newSub string
		start  int
		end    int
		result string
	}{
		{
			"U is a content management system",
			"UMONO",
			0,
			1,
			"UMONO is a content management system",
		},
		{
			"HELLO WORLD 0",
			"123",
			12,
			13,
			"HELLO WORLD 123",
		},
		{
			"HELLO WORLD 12345",
			"12345",
			12,
			17,
			"HELLO WORLD 12345",
		},
		{
			"HELLO great WORLD",
			"perfect",
			6,
			11,
			"HELLO perfect WORLD",
		},
		{
			"HELLO",
			"no-matter",
			-1,
			5,
			"HELLO",
		},
		{
			"HELLO",
			"no-matter",
			0,
			7,
			"HELLO",
		},
		{
			"HELLO",
			"no-matter",
			3,
			2,
			"HELLO",
		},
		// UTF-8
		{
			"öç öç öç öç şi şi",
			"perfect",
			6,
			8,
			"öç öç perfect öç şi şi",
		},
		{
			"öç öç öç öç şi şi",
			"perfect",
			0,
			2,
			"perfect öç öç öç şi şi",
		},
		{
			"öç öç öç öç şi şi",
			"perfect",
			15,
			17,
			"öç öç öç öç şi perfect",
		},
	} {
		result := ReplaceSubstring(scene.str, scene.newSub, scene.start, scene.end)
		require.Equal(s.T(), scene.result, result)
	}
}

func (s *StringsTestSuite) TestIsNumericScreamingSnakeCase() {
	for sI, scene := range []struct {
		str    string
		result bool
	}{
		{
			"UMONO_LANG",
			true,
		},
		{
			"UMONO__LANG",
			false,
		},
		{
			"UMONo",
			false,
		},
		{
			"_UMONO_LANG",
			false,
		},
		{
			" UMONO_LANG  ",
			false,
		},
		{
			"UMONO_123",
			true,
		},
		{
			"123_UMONO",
			true,
		},
		{
			"12345",
			true,
		},
	} {
		require.Equal(s.T(), scene.result, IsNumericScreamingSnakeCase(scene.str), "scene index: "+strconv.Itoa(sI))
	}
}

func TestStringsTestSuite(t *testing.T) {
	suite.Run(t, new(StringsTestSuite))
}
