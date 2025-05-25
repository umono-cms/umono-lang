package interfaces

type Component interface {
	Name() string
	Parameters() []Parameter
	RawContent() string
	NeedToConvert() bool
}
