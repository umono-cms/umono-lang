package umonolang

import (
	"fmt"

	"github.com/umono-cms/umono-lang/interfaces"
	ustrings "github.com/umono-cms/umono-lang/internal/utils/strings"
)

type call struct {
	component interfaces.Component
	start     int
	end       int
	args      []interfaces.Argument
}

func newCall(component interfaces.Component, start int, end int) *call {
	return &call{
		component: component,
		start:     start,
		end:       end,
	}
}

func (c *call) Component() interfaces.Component {
	return c.component
}

func (c *call) Start() int {
	return c.start
}

func (c *call) End() int {
	return c.end
}

func (c *call) Arguments() []interfaces.Argument {
	return c.args
}

func (c *call) ArgumentByName(name string) interfaces.Argument {
	for _, arg := range c.args {
		if arg.Name() == name {
			return arg
		}
	}
	return nil
}

func (c *call) fillParamsByRegex(content, regex string) {

	keyValueIndexes := ustrings.FindAllStringIndex(content, regex)
	params := readParams(content, keyValueIndexes)

	for _, compParam := range c.component.Parameters() {
		for _, param := range params {
			if compParam.Name() == param.Name() && compParam.Type() == param.Type() {
				c.args = append(c.args, newArg(param.Name(), param.Type(), param.Default()))
			}
		}
	}

	for _, compParam := range c.component.Parameters() {
		if filled := c.ArgumentByName(compParam.Name()); filled == nil {
			c.args = append(c.args, newArg(compParam.Name(), compParam.Type(), compParam.Default()))
		}
	}
}

type selector struct {
	regex      string
	paramRegex string
}

func readCalls(content string, comps []interfaces.Component) []*call {

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

	calls := []*call{}

	for _, slc := range selectors {
		for _, comp := range comps {
			indexes := ustrings.FindAllStringIndex(content, fmt.Sprintf(slc.regex, comp.Name()))
			for _, index := range indexes {
				if !alreadyRead(calledIndex, index[0], index[1]) {
					call := newCall(comp, index[0], index[1])
					call.fillParamsByRegex(string([]rune(content)[call.start:call.end]), slc.paramRegex)
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

func sortCallsByLinear(calls []*call) []*call {

	sorted := []*call{}
	clls := append([]*call{}, calls...)

	for {
		l, remainder := least(clls)
		if l == nil {
			break
		}

		sorted = append([]*call{l}, sorted...)
		clls = remainder
	}

	return sorted
}

func least(calls []*call) (*call, []*call) {

	if len(calls) == 0 {
		return nil, []*call{}
	}

	l := calls[0]
	start := calls[0].Start()

	for _, call := range calls {
		if start < call.Start() {
			l = call
			start = call.Start()
		}
	}

	remainder := []*call{}
	for _, call := range calls {
		if start != call.Start() {
			remainder = append(remainder, call)
		}
	}

	return l, remainder
}
