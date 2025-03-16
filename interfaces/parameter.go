package interfaces

type Parameter interface {
	Name() string
	Value() any
	SetValue(any)
}
