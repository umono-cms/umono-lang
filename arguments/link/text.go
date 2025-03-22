package link

type Text struct {
}

func (*Text) Name() string {
	return "text"
}

func (*Text) Type() string {
	return "string"
}

func (*Text) Default() any {
	return ""
}
