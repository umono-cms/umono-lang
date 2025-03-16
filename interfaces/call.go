package interfaces

type Call interface {
	Component() Component
	Start() int
	End() int
	Parameters() []Parameter
	ParameterByName(string) Parameter
}
