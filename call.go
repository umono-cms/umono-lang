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

func readCalls(content string, comps []interfaces.Component) []*Call {

	calls := []*Call{}

	for _, comp := range comps {
		indexes := ustrings.FindAllStringIndex(content, fmt.Sprintf(`%s`, comp.Name()))
		for _, index := range indexes {
			calls = append(calls, NewCall(comp, index[0], index[1]))
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
