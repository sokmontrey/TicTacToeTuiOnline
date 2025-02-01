package pageMsg

type ErrMsg struct {
	Value any
}

func NewErrMsg(value any) ErrMsg {
	return ErrMsg{value}
}
