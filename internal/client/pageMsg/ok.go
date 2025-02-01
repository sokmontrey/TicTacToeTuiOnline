package pageMsg

type OkMsg struct {
	Value any
}

func NewOkMsg(value any) OkMsg {
	return OkMsg{value}
}
