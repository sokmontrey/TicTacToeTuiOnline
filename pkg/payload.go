package pkg

import (
	"encoding/json"
	"github.com/eiannone/keyboard"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PayloadType byte

type MoveCode byte

type PositionUpdate struct {
	PlayerId int  `json:"playerId"`
	Position Vec2 `json:"position"`
}

type BoardUpdate struct {
	CellPos  Vec2 `json:"cellPos"`
	CellId   int  `json:"cellId"`
	NextTurn int  `json:"nextTurn"`
}

const (
	ServerErrPayload PayloadType = iota
	ServerOkPayload
	ServerJoinedIdPayload
	ServerPositionPayload
	ServerBoardUpdatePayload
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

func NewNonePayload() Payload {
	return NewPayload(NonePayload, nil)
}

func NewJoinedIdPayload(playerId int) Payload {
	return NewPayload(ServerJoinedIdPayload, playerId)
}

func NewPositionUpdatePayload(playerId int, position Vec2) Payload {
	return NewPayload(ServerPositionPayload, PositionUpdate{playerId, position})
}

func NewBoardUpdatePayload(cellPos Vec2, cellId int, nextTurn int) Payload {
	return NewPayload(ServerBoardUpdatePayload, BoardUpdate{
		cellPos,
		cellId,
		nextTurn,
	})
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
