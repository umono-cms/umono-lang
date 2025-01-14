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
