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

type PositionMsg struct {
	PlayerId int
	Position pkg.Vec2
}

type JoinedIdMsg struct {
	PlayerId int
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

func NewPositionMsg(playerId int, position pkg.Vec2) PositionMsg {
	return PositionMsg{playerId, position}
}

func NewJoinedIdMsg(playerId int) JoinedIdMsg {
	return JoinedIdMsg{playerId}
}
