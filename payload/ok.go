package payload

import "encoding/json"

type OkPayload struct {
	Value string `json:"value"`
}

func NewOkPayload(msg string) RawPayload {
	return NewPayload(ServerOkPayload, OkPayload{
		Value: msg,
	})
}

func (rp RawPayload) ToOkPayload() OkPayload {
	var okPayload OkPayload
	json.Unmarshal(rp.Data, &okPayload)
	return okPayload
}
