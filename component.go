package umonolang

import (
	"strings"

	"github.com/umono-cms/umono-lang/arguments"
	"github.com/umono-cms/umono-lang/components"
	"github.com/umono-cms/umono-lang/interfaces"
	ustrings "github.com/umono-cms/umono-lang/utils/strings"
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
	}

	argsKeyValueRaw := ustrings.FindAllString(raw, `([\w-]+)\s*=\s*"([^"]*)"|([\w-]+)\s*=\s*(true|false)`, "")

	args := []interfaces.Argument{}

	for _, keyValueRaw := range argsKeyValueRaw {
		ok, key, value := ustrings.SeparateKeyValue(keyValueRaw, `\s*=\s*`, `\s*"\s*|\s*"\s*`)
		if !ok {
			continue
		}

		typ := detectArgType(value)
		args = append(args, arguments.NewDynamicArg(key, typ, value))
	}

	rawContent := getRawContent(raw)

	return components.NewCustomWithArgs(name, strings.TrimSpace(rawContent), args)
}

func getCompName(raw string) string {
	names := ustrings.FindAllString(raw, `~\s+[A-Z0-9_]+(?:_[A-Z0-9]+)*\s*`, `(?s)^~\s*|\s*\n$`)
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func getRawContent(raw string) string {

	indexes := ustrings.FindAllStringIndex(raw, `([\w-]+)\s*=\s*"([^"]*)"|([\w-]+)\s*=\s*(true|false)`)
	if len(indexes) == 0 {
		indexes = ustrings.FindAllStringIndex(raw, `~\s+[A-Z0-9_]+(?:_[A-Z0-9]+)*\s*`)
		if len(indexes) == 0 {
			return raw
		}
	}

	last := indexes[len(indexes)-1]

	return string([]rune(raw)[last[1]:])
}

func detectArgType(value string) string {
	if value == "true" || value == "false" {
		return "bool"
	}
	return "string"
}
