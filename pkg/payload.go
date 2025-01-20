package pkg

import "encoding/json"

type PayloadType int
type KeyCode byte

const (
	ResErrPayloadType PayloadType = iota
	ResOkPayloadType
	ReqKeypressPayloadType
)

const (
	KeyCodeEsc KeyCode = iota
	KeyCodeConfirm
	KeyCodeUp
	KeyCodeDown
	KeyCodeLeft
	KeyCodeRight
)

type Payload struct {
	Type PayloadType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ResponsePayload struct {
	Type PayloadType `json:"type"`
	Data any         `json:"data"`
}

func NewOkResPayload(data any) Payload {
	rawData, _ := json.Marshal(data)
	return Payload{
		Type: ResOkPayloadType,
		Data: rawData,
	}
}

func NewErrResPayload(data any) Payload {
	rawData, _ := json.Marshal(data)
	return Payload{
		Type: ResErrPayloadType,
		Data: rawData,
	}
}

func NewKeypressReqPayload(key KeyCode) Payload {
	rawData, _ := json.Marshal(key)
	return Payload{
		Type: ReqKeypressPayloadType,
		Data: rawData,
	}
}
