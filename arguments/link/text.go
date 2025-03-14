package link

type Text struct {
	value any
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

func (t *Text) SetValue(val any) {
	t.value = val
}

func (t *Text) Value() any {
	return t.value
}
