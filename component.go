package umonolang

import (
	"github.com/umono-cms/umono-lang/components"
	"github.com/umono-cms/umono-lang/interfaces"
)

func builtInComps() []interfaces.Component {
	bcm := []interfaces.Component{}

	builtInComps := []interfaces.Component{
		&components.Link{},
	}

	for _, bc := range builtInComps {
		bcm = append(bcm, bc)
	}

	return bcm
}

func overrideComps(comps []interfaces.Component, dominantComps []interfaces.Component) []interfaces.Component {
	overridden := append([]interfaces.Component{}, comps...)

	for _, dc := range dominantComps {
		i, found := findCompByName(overridden, dc.Name())
		if found == nil {
			overridden = append(overridden, dc)
		} else {
			overridden[i] = dc
		}
	}

	return overridden
}

func findCompByName(comps []interfaces.Component, name string) (int, interfaces.Component) {
	for i, c := range comps {
		if c.Name() == name {
			return i, c
		}
	}
	return 0, nil
}
