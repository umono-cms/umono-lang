package umonolang

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/umono-cms/umono-lang/utils/mocks"
	"github.com/umono-cms/umono-lang/utils/test"
)

type UmonoLangTestSuite struct {
	suite.Suite
	umonoLang *UmonoLang
}

func (s *UmonoLangTestSuite) SetupTest() {

}

func (s *UmonoLangTestSuite) TestConvert() {

	converter := new(mocks.Converter)
	umonoLang := New(converter)

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

		require.Equal(s.T(), outputCont, umonoLang.Convert(inputCont), "input file name: "+input+".ul")
	}
}

func TestUmonoLangTestSuite(t *testing.T) {
	suite.Run(t, new(UmonoLangTestSuite))
}
