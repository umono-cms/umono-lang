package umonolang

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CallTestSuite struct {
	suite.Suite
}

func (s *CallTestSuite) TestAlreadyRead() {
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

func TestCallTestSuite(t *testing.T) {
	suite.Run(t, new(CallTestSuite))
}
