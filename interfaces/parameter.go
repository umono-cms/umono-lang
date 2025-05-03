package interfaces

type Parameter interface {
	Name() string
	Type() string
	Value() any
	SetValue(any)
	ValueAsString() string
}
