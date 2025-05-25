package link

type URL struct {
}

func (*URL) Name() string {
	return "url"
}

func (*URL) Type() string {
	return "string"
}

func (*URL) Default() any {
	return ""
}
