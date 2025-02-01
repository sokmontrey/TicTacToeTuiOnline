package payload

import (
	"encoding/json"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type PlayerPayload struct {
	PlayerId int      `json:"playerId"`
	Position pkg.Vec2 `json:"position"`
}

func NewPlayerPayload(playerId int, position pkg.Vec2) RawPayload {
	return NewPayload(ServerPlayerPayload, PlayerPayload{
		PlayerId: playerId,
		Position: position,
	})
}

func (rp RawPayload) ToPlayerPayload() PlayerPayload {
	var playerUpdatePayload PlayerPayload
	json.Unmarshal(rp.Data, &playerUpdatePayload)
	return playerUpdatePayload
}
