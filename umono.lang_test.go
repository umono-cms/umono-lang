package umonolang

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/umono-cms/umono-lang/utils/mocks"
)

type UmonoLangTestSuite struct {
	suite.Suite
	umonoLang *UmonoLang
}

func (s *UmonoLangTestSuite) SetupTest() {

}

func (s *UmonoLangTestSuite) TestConvertBasics() {

	converter := new(mocks.Converter)
	umonoLang := New(converter)

	for _, scene := range []struct {
		input  string
		output string
	}{
		{
			input:  "",
			output: "",
		},
		{
			input:  "input",
			output: "output",
		},
		{
			input:  "input\ninput",
			output: "outputoutput",
		},
		{
			input:  "{{COMPONENT}}\n~COMPONENT\ninput",
			output: "output",
		},
		{
			input:  "input{{COMPONENT}}\n~COMPONENT\ninput",
			output: "outputoutput",
		},
		{
			input:  "{{HEADER}} {{CONTENT}} \n~HEADER\ninput\n~CONTENT\ninput",
			output: "output output",
		},
		{
			input:  "{{HEADER}} {{CONTENT}} {{ FOOTER }} \n~HEADER\ninput\n~CONTENT\ninput\n~FOOTER\ninput",
			output: "output output output",
		},
	} {
		require.Equal(s.T(), scene.output, umonoLang.Convert(scene.input))
	}
}

func TestUmonoLangTestSuite(t *testing.T) {
	suite.Run(t, new(UmonoLangTestSuite))
}
