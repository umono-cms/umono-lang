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
	params    []interfaces.Parameter
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

func (c *Call) Parameters() []interfaces.Parameter {
	return c.params
}

func (c *Call) ParameterByName(name string) interfaces.Parameter {
	for _, p := range c.params {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

func (c *Call) fillArgsByRegex(content, regex, separator, trimRegex string) {

	keyValue := make(map[string]string)
	keyValuesRaw := ustrings.FindAllString(content, regex, "")

	for _, rawKeyVal := range keyValuesRaw {
		ok, key, value := ustrings.SeparateKeyValue(rawKeyVal, separator, trimRegex)
		if !ok {
			continue
		}

		keyValue[key] = value
	}

	for _, arg := range c.component.Arguments() {
		val, ok := keyValue[arg.Name()]
		if !ok {
			c.params = append(c.params, NewParam(arg.Name(), arg.Default()))
			continue
		}

		if arg.Type() == "string" {
			c.params = append(c.params, NewParam(arg.Name(), val))
			continue
		} else if arg.Type() == "bool" {
			valBool := true
			if val == "false" {
				valBool = false
			}
			c.params = append(c.params, NewParam(arg.Name(), valBool))
			continue
		}
	}

}

type selector struct {
	regex             string
	paramRegex        string
	keyValueSeparator string
	keyValueTrimRegex string
}

func readCalls(content string, comps []interfaces.Component) []*Call {

	// To prevent read substring regex match
	calledIndex := [][2]int{}

	selectors := []selector{
		selector{
			regex:             `\{\{\s*%s\s*([a-z0-9\-]+\s*=\s*.*)+\s*\}\}`,
			paramRegex:        `([\w-]+)\s*=\s*&quot;([^&]+)&quot;|([\w-]+)\s*=\s*(true|false)`,
			keyValueSeparator: `\s*=\s*`,
			keyValueTrimRegex: `\s*&quot;\s*`,
		},
		selector{
			regex:             `\{\{\s*%s\s*([a-z0-9\-]+\s*=\s*.*)+\s*\}\}`,
			paramRegex:        `([\w-]+)\s*=\s*"([^"]+)"|([\w-]+)\s*=\s*(true|false)`,
			keyValueSeparator: `\s*=\s*`,
			keyValueTrimRegex: `[\n\t\r\s]+`,
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
				if !alreadyRead(calledIndex, index[0], index[1]) {
					call := NewCall(comp, index[0], index[1])
					if slc.paramRegex != "" {
						call.fillArgsByRegex(string([]rune(content)[call.start:call.end]), slc.paramRegex, slc.keyValueSeparator, slc.keyValueTrimRegex)
					}
					calls = append(calls, call)
				}
			}
		}
	}

	return sortCallsByLinear(calls)
}

func alreadyRead(indexes [][2]int, start, end int) bool {
	return false
	for _, index := range indexes {
		if (start >= index[0] && start <= index[1]) || (end >= index[0] && end <= index[1]) {
			return true
		}
	}
	return false
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
