package payload

import (
	"encoding/json"
	"github.com/eiannone/keyboard"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
)

type MoveCode byte

const (
	MoveCodeNone MoveCode = iota
	MoveCodeConfirm
	MoveCodeUp
	MoveCodeDown
	MoveCodeLeft
	MoveCodeRight
)

func NewMoveCodePayload(moveCode MoveCode) RawPayload {
	return NewPayload(ClientMoveCodePayload, moveCode)
}

func (rp RawPayload) ToMoveCodePayload() MoveCode {
	var moveCode MoveCode
	json.Unmarshal(rp.Data, &moveCode)
	return moveCode
}

func KeyMsgToMoveCode(key pageMsg.KeyMsg) MoveCode {
	moveCode := keyboardToMoveCode(key.Key)
	if moveCode == MoveCodeNone {
		moveCode = charToMoveCode(key.Char)
	}
	return moveCode
}

func charToMoveCode(char rune) MoveCode {
	mapCharToMoveCode := map[rune]MoveCode{
		'w': MoveCodeUp,
		's': MoveCodeDown,
		'a': MoveCodeLeft,
		'd': MoveCodeRight,
		' ': MoveCodeConfirm,
	}
	moveCode, ok := mapCharToMoveCode[char]
	if ok {
		return moveCode
	}
	return MoveCodeNone
}

func keyboardToMoveCode(key keyboard.Key) MoveCode {
	mapKeyToMoveCode := map[keyboard.Key]MoveCode{
		keyboard.KeyArrowUp:    MoveCodeUp,
		keyboard.KeyArrowDown:  MoveCodeDown,
		keyboard.KeyArrowLeft:  MoveCodeLeft,
		keyboard.KeyArrowRight: MoveCodeRight,
		keyboard.KeyEnter:      MoveCodeConfirm,
		keyboard.KeySpace:      MoveCodeConfirm,
	}
	moveCode, ok := mapKeyToMoveCode[key]
	if ok {
		return moveCode
	}
	return MoveCodeNone
}
