package umonolang

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/umono-cms/umono-lang/components"
	"github.com/umono-cms/umono-lang/utils/mocks"
	"github.com/umono-cms/umono-lang/utils/test"
)

type UmonoLangTestSuite struct {
	suite.Suite
	umonoLang *UmonoLang
}

func (s *UmonoLangTestSuite) SetupTest() {
	s.umonoLang = New(new(mocks.Converter))
}

func (s *UmonoLangTestSuite) TestConvert() {
	inputFileReader := test.NewFileReader("test_assets/main/inputs", "ul")
	outputFileReader := test.NewFileReader("test_assets/main/outputs", "mock")

	inputDirReader := test.NewDirectoryReader("test_assets/main/inputs")
	inputs, err := inputDirReader.ReadWithoutExt()

	require.Nil(s.T(), err)

	for _, input := range inputs {

		inputCont, err := inputFileReader.Read(input, false)
		require.Nil(s.T(), err)

		outputCont, err := outputFileReader.Read(input, true)
		require.Nil(s.T(), err)

		require.Equal(s.T(), outputCont, s.umonoLang.Convert(inputCont), "input file name: "+input+".ul")
	}
}

func (s *UmonoLangTestSuite) TestSetGlobalComponentOK() {
	s.umonoLang.SetGlobalComponent("HELLO_WORLD", "hello!")
	hello, ok := s.umonoLang.globalCompMap["HELLO_WORLD"]

	require.True(s.T(), ok)
	require.Equal(s.T(), "hello!", hello.RawContent())
}

func (s *UmonoLangTestSuite) TestSetGlobalComponentSyntaxError() {
	err := s.umonoLang.SetGlobalComponent("HELLO WORLD", "hello!")
	require.NotNil(s.T(), err)
	require.True(s.T(), strings.HasPrefix(err.Error(), "SYNTAX_ERROR"))
}

func (s *UmonoLangTestSuite) TestRemoveGlobalComponentOK() {
	s.umonoLang.globalCompMap["HELLO_WORLD"] = components.NewCustom("HELO_WORLD", "hello!")
	err := s.umonoLang.RemoveGlobalComponent("HELLO_WORLD")
	require.Nil(s.T(), err)

	_, ok := s.umonoLang.globalCompMap["HELLO_WORLD"]
	require.False(s.T(), ok)
}

func (s *UmonoLangTestSuite) TestRemoveGlobalComponentSyntaxError() {
	err := s.umonoLang.RemoveGlobalComponent("hello world")
	require.NotNil(s.T(), err)
	require.True(s.T(), strings.HasPrefix(err.Error(), "SYNTAX_ERROR"))
}

func (s *UmonoLangTestSuite) TestRemoveGlobalComponentNotFound() {
	err := s.umonoLang.RemoveGlobalComponent("HELLO_WORLD")
	require.NotNil(s.T(), err)
	require.True(s.T(), strings.HasPrefix(err.Error(), "NOT_FOUND"))
}

func TestUmonoLangTestSuite(t *testing.T) {
	suite.Run(t, new(UmonoLangTestSuite))
}
