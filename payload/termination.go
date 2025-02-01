package payload

import (
	"encoding/json"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type TerminationPayload struct {
	ConnectedCells []pkg.Vec2 `json:"connectedCells"`
	WinnerId       int        `json:"winnerId"`
}

func NewTerminationPayload(winnerId int, connectedCells map[pkg.Vec2]struct{}) RawPayload {
	connectedCellsArr := make([]pkg.Vec2, len(connectedCells))
	for v := range connectedCells {
		connectedCellsArr = append(connectedCellsArr, v)
	}
	return NewPayload(ServerTerminationPayload, TerminationPayload{
		WinnerId:       winnerId,
		ConnectedCells: connectedCellsArr,
	})
}

func (rp RawPayload) ToTerminationPayload() TerminationPayload {
	var terminationUpdate TerminationPayload
	json.Unmarshal(rp.Data, &terminationUpdate)
	return terminationUpdate
}

func (t TerminationPayload) GetConnectedCellsMap() map[pkg.Vec2]struct{} {
	connectedCells := make(map[pkg.Vec2]struct{})
	for _, v := range t.ConnectedCells {
		connectedCells[v] = struct{}{}
	}
	return connectedCells
}
