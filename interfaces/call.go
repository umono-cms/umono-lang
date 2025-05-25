package interfaces

type Call interface {
	Component() Component
	Start() int
	End() int
	Arguments() []Argument
	ArgumentByName(string) Argument
}
