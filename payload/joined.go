package payload

import "encoding/json"

type JoinedPayload struct {
	PlayerId int `json:"playerId"`
}

func NewJoinedPayload(playerId int) RawPayload {
	return NewPayload(ServerJoinedPayload, JoinedPayload{
		PlayerId: playerId,
	})
}

func (rp RawPayload) ToJoinedPayload() JoinedPayload {
	var joinedUpdate JoinedPayload
	json.Unmarshal(rp.Data, &joinedUpdate)
	return joinedUpdate
}
