package pkg

import "encoding/json"

type PayloadType int
type KeyCode byte

const (
	ServerErrPayloadType PayloadType = iota
	ServerOkPayloadType
	ClientKeypressPayloadType
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

type ServerPayload struct {
	Type PayloadType `json:"type"`
	Data any         `json:"data"`
}

func NewOkServerPayload(data any) Payload {
	rawData, _ := json.Marshal(data)
	return Payload{
		Type: ServerOkPayloadType,
		Data: rawData,
	}
}

func NewErrServerPayload(data any) Payload {
	rawData, _ := json.Marshal(data)
	return Payload{
		Type: ServerErrPayloadType,
		Data: rawData,
	}
}

func NewKeypressClientPayload(key KeyCode) Payload {
	rawData, _ := json.Marshal(key)
	return Payload{
		Type: ClientKeypressPayloadType,
		Data: rawData,
	}
}
