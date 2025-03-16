package interfaces

type Converter interface {
	Convert(string) string
	ConvertBuiltInComp(Call) string
}
