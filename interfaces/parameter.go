package interfaces

type Parameter interface {
	Name() string
	Type() string
	Default() any
}
