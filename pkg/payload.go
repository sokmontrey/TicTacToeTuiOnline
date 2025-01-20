package pkg

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
	Type PayloadType `json:"type"`
	Data any         `json:"data"`
}

func NewOkResPayload(data any) Payload {
	return Payload{
		Type: ResOkPayloadType,
		Data: data,
	}
}

func NewErrResPayload(data any) Payload {
	return Payload{
		Type: ResErrPayloadType,
		Data: data,
	}
}

func NewKeypressReqPayload(key KeyCode) Payload {
	return Payload{
		Type: ReqKeypressPayloadType,
		Data: key,
	}
}
