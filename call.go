package umonolang

import (
	"fmt"

	"github.com/umono-cms/umono-lang/interfaces"
	ustrings "github.com/umono-cms/umono-lang/internal/utils/strings"
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

func (c *Call) fillArgsByRegex(content, regex string) {

	keyValueIndexes := ustrings.FindAllStringIndex(content, regex)
	args := readArgs(content, keyValueIndexes)

	for _, compArg := range c.component.Arguments() {
		for _, arg := range args {
			if compArg.Name() == arg.Name() && compArg.Type() == arg.Type() {
				c.params = append(c.params, NewParam(arg.Name(), arg.Type(), arg.Default()))
			}
		}
	}

	for _, compArg := range c.component.Arguments() {
		if filled := c.ParameterByName(compArg.Name()); filled == nil {
			c.params = append(c.params, NewParam(compArg.Name(), compArg.Type(), compArg.Default()))
		}
	}
}

type selector struct {
	regex      string
	paramRegex string
}

func readCalls(content string, comps []interfaces.Component) []*Call {

	calledIndex := [][2]int{}

	selectors := []selector{
		selector{
			regex:      `\{\{\s*%s(?:\s+(?:[\w-]+\s*=\s*&quot;[^&]+&quot;|[\w-]+\s*=\s*(?:true|false)))*\s*\}\}`,
			paramRegex: `([\w-]+)\s*=\s*`,
		},
		selector{
			regex:      `\{\{\s*%s(?:\s+[\w-]+\s*=\s*.*\s*)+\s*\}\}`,
			paramRegex: `([\w-]+)\s*=\s*`,
		},
		selector{
			regex: `\{\{\s*%s\s*\}\}`,
		},
		selector{
			regex: `%s`,
		},
	}

	calls := []*Call{}

	for _, slc := range selectors {
		for _, comp := range comps {
			indexes := ustrings.FindAllStringIndex(content, fmt.Sprintf(slc.regex, comp.Name()))
			for _, index := range indexes {
				if !alreadyRead(calledIndex, index[0], index[1]) {
					call := NewCall(comp, index[0], index[1])
					call.fillArgsByRegex(string([]rune(content)[call.start:call.end]), slc.paramRegex)
					calls = append(calls, call)
					calledIndex = append(calledIndex, [2]int{index[0], index[1]})
				}
			}
		}
	}

	return sortCallsByLinear(calls)
}

func alreadyRead(indexes [][2]int, start, end int) bool {
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
