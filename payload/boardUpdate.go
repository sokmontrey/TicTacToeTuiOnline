package payload

import (
	"encoding/json"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type BoardUpdatePayload struct {
	Cell     CellPayload `json:"cell"`
	NextTurn int         `json:"nextTurn"`
}

func NewBoardUpdatePayload(cellPos pkg.Vec2, cellId int, nextTurn int) RawPayload {
	return NewPayload(ServerBoardUpdatePayload, BoardUpdatePayload{
		Cell: CellPayload{
			CellPos: cellPos,
			CellId:  cellId,
		},
		NextTurn: nextTurn,
	})
}

func (rp RawPayload) ToBoardUpdatePayload() BoardUpdatePayload {
	var boardUpdatePayload BoardUpdatePayload
	json.Unmarshal(rp.Data, &boardUpdatePayload)
	return boardUpdatePayload
}
