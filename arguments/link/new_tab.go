package link

type NewTab struct {
	value any
}

func (*NewTab) Name() string {
	return "new-tab"
}

func (*NewTab) Type() string {
	return "bool"
}

func (*NewTab) Default() any {
	return false
}

func (nt *NewTab) SetValue(val any) {
	nt.value = val
}

func (nt *NewTab) Value() any {
	return nt.value
}
