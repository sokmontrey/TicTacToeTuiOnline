package payload

import (
	"encoding/json"
	"github.com/eiannone/keyboard"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type PayloadType byte

type MoveCode byte

type PlayerUpdate struct {
	PlayerId int      `json:"playerId"`
	Position pkg.Vec2 `json:"position"`
}

type CellUpdate struct {
	CellPos pkg.Vec2 `json:"cellPos"`
	CellId  int      `json:"cellId"`
}

type BoardUpdate struct {
	Cell     CellUpdate `json:"cell"`
	NextTurn int        `json:"nextTurn"`
}

type SyncUpdate struct {
	PlayerPositions []PlayerUpdate `json:"playerPositions"`
	CellPositions   []CellUpdate   `json:"cellPositions"`
	CurrentTurn     int            `json:"currentTurn"`
	CurrentPlayerId int            `json:"currentPlayerId"`
}

type TerminationUpdate struct {
	ConnectedCells []pkg.Vec2 `json:"connectedCells"`
	WinnerId       int        `json:"winnerId"`
}

const (
	ServerErrPayload PayloadType = iota
	ServerOkPayload
	ServerSyncPayload
	ServerPositionPayload
	ServerJoinedPayload
	ServerTerminationPayload
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

func NewSyncPayload(playerPositions []PlayerUpdate,
	cellPositions []CellUpdate,
	currentTurn int,
	currentPlayerId int,
) RawPayload {
	return NewPayload(ServerSyncPayload, SyncUpdate{
		PlayerPositions: playerPositions,
		CellPositions:   cellPositions,
		CurrentTurn:     currentTurn,
		CurrentPlayerId: currentPlayerId,
	})
}

func NewPositionUpdatePayload(playerId int, position pkg.Vec2) RawPayload {
	return NewPayload(ServerPositionPayload, PlayerUpdate{playerId, position})
}

func NewBoardUpdatePayload(cellPos pkg.Vec2, cellId int, nextTurn int) RawPayload {
	return NewPayload(ServerBoardUpdatePayload, BoardUpdate{
		Cell: CellUpdate{
			CellPos: cellPos,
			CellId:  cellId,
		},
		NextTurn: nextTurn,
	})
}

func NewTerminationPayload(winnerId int, connectedCells map[pkg.Vec2]struct{}) RawPayload {
	connectedCellsArr := make([]pkg.Vec2, len(connectedCells))
	for v := range connectedCells {
		connectedCellsArr = append(connectedCellsArr, v)
	}
	return NewPayload(ServerTerminationPayload, TerminationUpdate{
		WinnerId:       winnerId,
		ConnectedCells: connectedCellsArr,
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
