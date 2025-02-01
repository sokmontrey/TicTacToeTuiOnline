package payload

import "encoding/json"

type ErrPayload struct {
	Value string `json:"value"`
}

func NewErrPayload(msg string) RawPayload {
	return NewPayload(ServerErrPayload, ErrPayload{
		Value: msg,
	})
}

func (rp RawPayload) ToErrPayload() ErrPayload {
	var okPayload ErrPayload
	json.Unmarshal(rp.Data, &okPayload)
	return okPayload
}
