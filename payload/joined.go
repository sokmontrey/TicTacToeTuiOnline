package payload

import "encoding/json"

type JoinedUpdate struct {
	PlayerId int `json:"playerId"`
}

func NewJoinedUpdatePayload(playerId int) RawPayload {
	return NewPayload(ServerJoinedPayload, JoinedUpdate{
		PlayerId: playerId,
	})
}

func (rp RawPayload) ToJoinedUpdate() JoinedUpdate {
	var joinedUpdate JoinedUpdate
	json.Unmarshal(rp.Data, &joinedUpdate)
	return joinedUpdate
}
