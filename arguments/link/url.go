package link

type URL struct {
	value any
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

func (u *URL) SetValue(val any) {
	u.value = val
}

func (u *URL) Value() any {
	return u.value
}
