package interfaces

type Argument interface {
	Name() string
	Type() string
	Default() any
}
