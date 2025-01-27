package pageMsg

import (
	"github.com/eiannone/keyboard"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type PageMsg any

type KeyMsg struct {
	Char rune
	Key  keyboard.Key
}

func (k KeyMsg) ToMoveCode() pkg.MoveCode {
	moveCode := pkg.KeyToMoveCode(k.Key)
	if moveCode == pkg.MoveCodeNone {
		moveCode = pkg.CharToMoveCode(k.Char)
	}
	return moveCode
}

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
