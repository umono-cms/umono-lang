package umonolang

import (
	"fmt"

	"github.com/umono-cms/umono-lang/interfaces"
	ustrings "github.com/umono-cms/umono-lang/utils/strings"
)

type Call struct {
	component interfaces.Component
	start     int
	end       int
}

func NewCall(component interfaces.Component, start int, end int) *Call {
	return &Call{
		component: component,
		start:     start,
		end:       end,
	}
}

func (c *Call) Component() interfaces.Component {
	return c.component
}

func (c *Call) Start() int {
	return c.start
}

func (c *Call) End() int {
	return c.end
}

func (c *Call) fillArgsByRegex(content, regex, separator string) {

	keyValue := make(map[string]string)
	keyValuesRaw := ustrings.FindAllString(content, regex, "")

	for _, rawKeyVal := range keyValuesRaw {
		ok, key, value := ustrings.SeparateKeyValue(rawKeyVal, separator)
		if !ok {
			continue
		}

		keyValue[key] = value
	}

	for _, arg := range c.component.Arguments() {
		val, ok := keyValue[arg.Name()]
		if !ok {
			arg.SetValue(arg.Default())
			continue
		}

		if arg.Type() == "string" {
			arg.SetValue(val)
		} else if arg.Type() == "bool" {
			if val == "false" {
				arg.SetValue(false)
			} else if val == "true" {
				arg.SetValue(false)
			}
		}
	}

}

type selector struct {
	regex             string
	paramRegex        string
	keyValueSeparator string
}

func readCalls(content string, comps []interfaces.Component) []*Call {

	selectors := []selector{
		selector{
			regex:             `\{\{\s*%s\s*([a-z0-9\-]+\s*=\s*.*)+\s*\}\}`,
			paramRegex:        `([\w-]+)\s*=\s*"([^"]+)"|([\w-]+)\s*=\s*(true|false)`,
			keyValueSeparator: `\s*=\s*`,
		},
		selector{
			regex:      `\{\{\s*%s\s*\}\}`,
			paramRegex: "",
		},
		selector{
			regex:      `%s`,
			paramRegex: "",
		},
	}

	calls := []*Call{}

	for _, slc := range selectors {
		for _, comp := range comps {
			indexes := ustrings.FindAllStringIndex(content, fmt.Sprintf(slc.regex, comp.Name()))
			for _, index := range indexes {
				call := NewCall(comp, index[0], index[1])
				if slc.paramRegex != "" {
					call.fillArgsByRegex(string([]rune(content)[call.start:call.end]), slc.paramRegex, slc.keyValueSeparator)
				}
				calls = append(calls, call)
			}
		}
	}

	return sortCallsByLinear(calls)
}

func sortCallsByLinear(calls []*Call) []*Call {

	sorted := []*Call{}
	clls := append([]*Call{}, calls...)

	for {
		l, remainder := least(clls)
		if l == nil {
			break
		}

		sorted = append([]*Call{l}, sorted...)
		clls = remainder
	}

	return sorted
}

func least(calls []*Call) (*Call, []*Call) {

	if len(calls) == 0 {
		return nil, []*Call{}
	}

	l := calls[0]
	start := calls[0].Start()

	for _, call := range calls {
		if start < call.Start() {
			l = call
			start = call.Start()
		}
	}

	remainder := []*Call{}
	for _, call := range calls {
		if start != call.Start() {
			remainder = append(remainder, call)
		}
	}

	return l, remainder
}
