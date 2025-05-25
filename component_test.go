package umonolang

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/umono-cms/umono-lang/interfaces"
	"github.com/umono-cms/umono-lang/internal/components"
)

type ComponentTestSuite struct {
	suite.Suite
}

func (suite *ComponentTestSuite) TestOverrideComps() {
	comps := []interfaces.Component{
		components.NewCustom("A", "a"),
		components.NewCustom("B", "b"),
	}

	otherComps := []interfaces.Component{
		components.NewCustom("C", "a"),
		components.NewCustom("D", "b"),
	}

	overridden := overrideComps(comps, otherComps)

	require.Equal(suite.T(), int(2), len(comps))
	require.Equal(suite.T(), int(2), len(otherComps))
	require.Equal(suite.T(), int(4), len(overridden))
}

func TestComponentTestSuite(t *testing.T) {
	suite.Run(t, new(ComponentTestSuite))
}
