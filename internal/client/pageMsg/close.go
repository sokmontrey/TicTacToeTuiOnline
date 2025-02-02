package pageMsg

type CloseMsg struct {
	Value string
}

func NewCloseMsg(value string) CloseMsg {
	return CloseMsg{value}
}
