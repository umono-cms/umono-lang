package umonolang

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/umono-cms/umono-lang/interfaces"
	"github.com/umono-cms/umono-lang/internal/arguments"
	"github.com/umono-cms/umono-lang/internal/components"
)

type callTestSuite struct {
	suite.Suite
}

func (s *callTestSuite) TestAlreadyRead() {
	for sI, scene := range []struct {
		indexes [][2]int
		start   int
		end     int
		result  bool
	}{
		{
			[][2]int{
				[2]int{30, 40},
			},
			34,
			45,
			true,
		},
		{
			[][2]int{
				[2]int{30, 40},
				[2]int{5, 15},
			},
			1,
			3,
			false,
		},
		{
			[][2]int{
				[2]int{30, 40},
				[2]int{90, 120},
			},
			125,
			154,
			false,
		},
		{
			[][2]int{
				[2]int{30, 40},
				[2]int{90, 120},
			},
			50,
			60,
			false,
		},
	} {
		require.Equal(s.T(), scene.result, alreadyRead(scene.indexes, scene.start, scene.end), "scene index: "+strconv.Itoa(sI))
	}
}

func (s *callTestSuite) TestReadCalls() {
	for sI, scene := range []struct {
		content string
		comps   []interfaces.Component
		results []*call
	}{
		{
			"{{ ABC }}",
			[]interfaces.Component{
				components.NewCustom("ABC", "no-matter"),
			},
			[]*call{
				&call{components.NewCustom("ABC", "no-matter"), 0, 9, []interfaces.Parameter{}},
			},
		},
		{
			"{{ ABC param-1=\"val-1\" param-2=\"val-2\" }}",
			[]interfaces.Component{
				components.NewCustomWithArgs("ABC", "no-matter", []interfaces.Argument{
					arguments.NewDynamicArg("param-1", "string", "val-1"),
					arguments.NewDynamicArg("param-2", "string", "val-2"),
				}),
			},
			[]*call{
				&call{components.NewCustomWithArgs("ABC", "no-matter", []interfaces.Argument{
					arguments.NewDynamicArg("param-1", "string", ""),
					arguments.NewDynamicArg("param-2", "string", ""),
				}), 0, 41, []interfaces.Parameter{
					newParam("param-1", "string", "val-1"),
					newParam("param-2", "string", "val-2"),
				}},
			},
		},
		{
			"ABC XYZ",
			[]interfaces.Component{
				components.NewCustom("ABC", "no-matter"),
				components.NewCustom("XYZ", "no-matter"),
			},
			[]*call{
				&call{components.NewCustom("ABC", "no-matter"), 0, 3, []interfaces.Parameter{}},
				&call{components.NewCustom("XYZ", "no-matter"), 4, 7, []interfaces.Parameter{}},
			},
		},
		{
			"XYZ",
			[]interfaces.Component{
				components.NewCustomWithArgs("XYZ", "no-matter", []interfaces.Argument{
					arguments.NewDynamicArg("param", "string", "default-value"),
				}),
			},
			[]*call{
				&call{components.NewCustomWithArgs("XYZ", "no-matter", []interfaces.Argument{
					arguments.NewDynamicArg("param", "string", "default-value"),
				}), 0, 3, []interfaces.Parameter{
					newParam("param", "string", "default-value"),
				}},
			},
		},
		{
			"{{ LINK url=\"https://umono.io\" text=\"click me!\" new-tab=true }}",
			[]interfaces.Component{
				components.NewCustomWithArgs("LINK", "no-matter", []interfaces.Argument{
					arguments.NewDynamicArg("url", "string", ""),
					arguments.NewDynamicArg("text", "string", ""),
					arguments.NewDynamicArg("new-tab", "bool", false),
				}),
			},
			[]*call{
				&call{
					components.NewCustomWithArgs("LINK", "no-matter", []interfaces.Argument{
						arguments.NewDynamicArg("url", "string", ""),
						arguments.NewDynamicArg("text", "string", ""),
						arguments.NewDynamicArg("new-tab", "bool", false),
					}), 0, 63, []interfaces.Parameter{
						newParam("url", "string", "https://umono.io"),
						newParam("text", "string", "click me!"),
						newParam("new-tab", "bool", true),
					},
				},
			},
		},
	} {

		calls := readCalls(scene.content, scene.comps)

		errMsg := "scene index: " + strconv.Itoa(sI)

		for i := 0; i < len(scene.results); i++ {

			sCall := scene.results[i]
			call := calls[i]

			errMsg := errMsg + ", call index: " + strconv.Itoa(i)

			require.Equal(s.T(), sCall.Component().Name(), call.Component().Name(), errMsg)
			require.Equal(s.T(), sCall.Start(), call.Start(), errMsg)
			require.Equal(s.T(), sCall.End(), call.End(), errMsg)

			sParams := sCall.Parameters()
			callParams := call.Parameters()

			for j := 0; j < len(sParams); j++ {
				errMsg := errMsg + ", param index: " + strconv.Itoa(j)
				require.Equal(s.T(), sParams[j].Name(), callParams[j].Name(), errMsg)
				if sParams[j].Type() == "string" {
					require.Equal(s.T(), sParams[j].Value().(string), callParams[j].Value().(string), errMsg)
				} else if sParams[j].Type() == "bool" {
					require.Equal(s.T(), sParams[j].Value().(bool), callParams[j].Value().(bool), errMsg)
				}
			}
		}
	}
}

func TestCallTestSuite(t *testing.T) {
	suite.Run(t, new(callTestSuite))
}
