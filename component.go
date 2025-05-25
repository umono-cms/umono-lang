package umonolang

import (
	"strings"

	"github.com/umono-cms/umono-lang/interfaces"
	"github.com/umono-cms/umono-lang/internal/components"
	ustrings "github.com/umono-cms/umono-lang/internal/utils/strings"
)

func builtInComps() []interfaces.Component {
	bcm := []interfaces.Component{}

	builtInComps := []interfaces.Component{
		&components.Link{},
		&components.S404{},
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

func readLocalComps(localCompsRaw string) []interfaces.Component {

	localCompIndexes := ustrings.IndexesByRegex(localCompsRaw, `\n~\s+[A-Z0-9_]+(?:_[A-Z0-9]+)*\s*\n`)

	comps := []interfaces.Component{}

	for i := 0; i < len(localCompIndexes); i++ {
		var raw string
		if i == len(localCompIndexes)-1 {
			raw = localCompsRaw[localCompIndexes[i]:]
		} else {
			raw = localCompsRaw[localCompIndexes[i]:localCompIndexes[i+1]]
		}
		comps = append(comps, readComp("", strings.TrimSpace(raw)))
	}

	return comps
}

func readComp(compName, raw string) interfaces.Component {

	name := compName
	if compName == "" {
		name = getCompName(raw)
		raw = getRawWithoutCompName(raw)
	}

	keyValueIndexes := ustrings.FindAllStringIndex(raw, `([\w-]+)\s*=\s*`)
	args := readArgs(raw, keyValueIndexes)

	return components.NewCustomWithArgs(name, getContentBody(raw), args)
}

func getCompName(raw string) string {
	names := ustrings.FindAllString(raw, `~\s+[A-Z0-9_]+(?:_[A-Z0-9]+)*\s*`, `(?s)^~\s*|\s*\n$`)
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func getRawWithoutCompName(raw string) string {
	indexes := ustrings.FindAllStringIndex(raw, `~\s+[A-Z0-9_]+(?:_[A-Z0-9]+)*\s*`)
	if len(indexes) == 0 {
		return raw
	}

	return raw[indexes[0][1]:]
}

func getContentBody(raw string) string {

	indexes := ustrings.FindAllStringIndex(raw, `([\w-]+)\s*=\s*`)

	if len(indexes) == 0 {
		return strings.TrimSpace(raw)
	}

	last := indexes[len(indexes)-1]

	newLineIndex := strings.Index(raw[last[1]:], "\n")

	if newLineIndex == -1 {
		return strings.TrimSpace(raw)
	}

	return strings.TrimSpace(string(raw[last[1]+newLineIndex+1:]))
}
