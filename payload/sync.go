package payload

import (
	"encoding/json"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/game"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type SyncPayload struct {
	PlayerPositions []PlayerPayload `json:"playerPositions"`
	CellPositions   []CellPayload   `json:"cellPositions"`
	CurrentTurn     int             `json:"currentTurn"`
	CurrentPlayerId int             `json:"currentPlayerId"`
}

func NewSyncPayload(
	players map[int]*game.Player,
	cells map[pkg.Vec2]int,
	currentTurn int,
	currentPlayerId int,
) RawPayload {
	playerArr := make([]PlayerPayload, 0)
	cellArr := make([]CellPayload, 0)

	for id, player := range players {
		playerArr = append(playerArr, PlayerPayload{PlayerId: id, Position: player.Position})
	}
	for pos, cellId := range cells {
		cellArr = append(cellArr, CellPayload{CellId: cellId, CellPos: pos})
	}

	return NewPayload(ServerSyncPayload, SyncPayload{
		PlayerPositions: playerArr,
		CellPositions:   cellArr,
		CurrentTurn:     currentTurn,
		CurrentPlayerId: currentPlayerId,
	})
}

func (rp RawPayload) ToSyncUpdatePayload() SyncPayload {
	var syncUpdatePayload SyncPayload
	json.Unmarshal(rp.Data, &syncUpdatePayload)
	return syncUpdatePayload
}
