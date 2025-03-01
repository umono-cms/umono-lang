package interfaces

type Component interface {
	Name() string
	Arguments() []Argument
	RawContent() string
	PutAfterConvert() bool
}
