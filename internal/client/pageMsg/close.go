package pageMsg

type CloseMsg struct {
	Value any
}

func NewCloseMsg(value any) CloseMsg {
	return CloseMsg{value}
}
