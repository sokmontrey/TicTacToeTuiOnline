package payload

import (
	"encoding/json"
	"github.com/eiannone/keyboard"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PayloadType byte

type MoveCode byte

const (
	ServerErrPayload PayloadType = iota
	ServerOkPayload
	ServerJoinedPayload
	ServerPlayerUpdatePayload
	ServerBoardUpdatePayload
	ServerSyncPayload
	ServerTerminationPayload
	ClientMovePayload
	NonePayload
)

const (
	MoveCodeNone MoveCode = iota
	MoveCodeConfirm
	MoveCodeUp
	MoveCodeDown
	MoveCodeLeft
	MoveCodeRight
)

type RawPayload struct {
	Type PayloadType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (rp RawPayload) WsSend(conn *websocket.Conn) error {
	return conn.WriteJSON(rp)
}

func (rp RawPayload) HttpSend(statusCode int, c *gin.Context) error {
	c.JSON(statusCode, rp)
	return nil
}

func NewPayload(payloadType PayloadType, data any) RawPayload {
	rawData, _ := json.Marshal(data)
	return RawPayload{
		Type: payloadType,
		Data: rawData,
	}
}

func NewNonePayload() RawPayload {
	return NewPayload(NonePayload, nil)
}

func CharToMoveCode(char rune) MoveCode {
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

func KeyToMoveCode(key keyboard.Key) MoveCode {
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
