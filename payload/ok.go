package payload

import "encoding/json"

type OkPayload struct {
	Value string `json:"value"`
}

func NewOkPayload(value string) RawPayload {
	return NewPayload(ServerOkPayload, OkPayload{value})
}

func (rp RawPayload) ToOkPayload() OkPayload {
	var okPayload OkPayload
	json.Unmarshal(rp.Data, &okPayload)
	return okPayload
}
