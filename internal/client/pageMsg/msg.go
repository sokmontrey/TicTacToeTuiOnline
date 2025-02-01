package pageMsg

import (
	"github.com/eiannone/keyboard"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
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

func (k KeyMsg) ToMoveCode() payload.MoveCode {
	moveCode := payload.KeyToMoveCode(k.Key)
	if moveCode == payload.MoveCodeNone {
		moveCode = payload.CharToMoveCode(k.Char)
	}
	return moveCode
}
