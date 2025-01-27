package pkg

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
	ClientMovePayload
)

const (
	MoveCodeNone MoveCode = iota
	MoveCodeConfirm
	MoveCodeUp
	MoveCodeDown
	MoveCodeLeft
	MoveCodeRight
)

type Payload struct {
	Type PayloadType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (p Payload) WsSend(conn *websocket.Conn) error {
	return conn.WriteJSON(p)
}

func (p Payload) HttpSend(statusCode int, c *gin.Context) error {
	c.JSON(statusCode, p)
	return nil
}

func NewPayload(payloadType PayloadType, data any) Payload {
	rawData, _ := json.Marshal(data)
	return Payload{
		Type: payloadType,
		Data: rawData,
	}
}

func CharToMoveCode(char rune) MoveCode {
	mapCharToMoveCode := map[rune]MoveCode{
		'w': MoveCodeUp,
		's': MoveCodeDown,
		'a': MoveCodeLeft,
		'd': MoveCodeRight,
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
	}
	moveCode, ok := mapKeyToMoveCode[key]
	if ok {
		return moveCode
	}
	return MoveCodeNone
}
