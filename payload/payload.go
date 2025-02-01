package payload

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PayloadType byte

const (
	ServerErrPayload PayloadType = iota
	ServerOkPayload
	ServerJoinedPayload
	ServerPlayerPayload
	ServerBoardUpdatePayload
	ServerSyncPayload
	ServerTerminationPayload
	ClientMoveCodePayload
	NonePayload
)

type RawPayload struct {
	Type PayloadType     `json:"type"`
	Data json.RawMessage `json:"data"`
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

func (rp RawPayload) WsSend(conn *websocket.Conn) error {
	return conn.WriteJSON(rp)
}

func (rp RawPayload) HttpSend(statusCode int, c *gin.Context) error {
	c.JSON(statusCode, rp)
	return nil
}
