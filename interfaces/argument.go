package interfaces

type Argument interface {
	Name() string
	Type() string
	Value() any
	SetValue(any)
	ValueAsString() string
}
