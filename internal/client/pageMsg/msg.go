package pageMsg

import (
	"github.com/eiannone/keyboard"
)

type PageMsg any

type OkMsg struct {
	Data any
}

type ErrMsg struct {
	Data any
}

func NewKeyMsg(char rune, key keyboard.Key) KeyMsg {
	return KeyMsg{char, key}
}

func NewOkMsg(data any) OkMsg {
	return OkMsg{data}
}

func NewErrMsg(data any) ErrMsg {
	return ErrMsg{data}
}

type KeyMsg struct {
	Char rune
	Key  keyboard.Key
}
