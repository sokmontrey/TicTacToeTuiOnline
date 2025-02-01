package payload

import (
	"encoding/json"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type CellPayload struct {
	CellPos pkg.Vec2 `json:"cellPos"`
	CellId  int      `json:"cellId"`
}

func NewCellPayload(cellPos pkg.Vec2, cellId int) RawPayload {
	return NewPayload(ServerBoardUpdatePayload, CellPayload{
		CellPos: cellPos,
		CellId:  cellId,
	})
}

func (rp RawPayload) ToCellPayload() CellPayload {
	var cellPayload CellPayload
	json.Unmarshal(rp.Data, &cellPayload)
	return cellPayload
}
