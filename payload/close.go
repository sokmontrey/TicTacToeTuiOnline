package payload

import "encoding/json"

type ClosePayload struct {
	Msg string `json:"msg"`
}

func NewClosePayload(msg string) RawPayload {
	return NewPayload(ServerClosePayload, ClosePayload{msg})
}

func (rp RawPayload) ToClosePayload() ClosePayload {
	var closePayload ClosePayload
	json.Unmarshal(rp.Data, &closePayload)
	return closePayload
}
