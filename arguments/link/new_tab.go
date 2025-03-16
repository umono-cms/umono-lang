package link

type NewTab struct {
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
