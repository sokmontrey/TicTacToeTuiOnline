package pkg

import (
	"encoding/json"
	"github.com/eiannone/keyboard"
)

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
	KeyCodeNone
)

type Payload struct {
	Type PayloadType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ServerPayload struct {
	Type PayloadType `json:"type"`
	Data any         `json:"data"`
}

type ClientPayload struct {
	Type PayloadType `json:"type"`
	Data KeyCode     `json:"data"`
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

func CharToKeyCode(char rune) KeyCode {
	switch char {
	case 'w':
		return KeyCodeUp
	case 's':
		return KeyCodeDown
	case 'a':
		return KeyCodeLeft
	case 'd':
		return KeyCodeRight
	case ' ':
		return KeyCodeConfirm
	}
	return KeyCodeNone
}

func KeyPressToKeyCode(key keyboard.Key) KeyCode {
	switch key {
	case keyboard.KeyArrowUp:
		return KeyCodeUp
	case keyboard.KeyArrowDown:
		return KeyCodeDown
	case keyboard.KeyArrowLeft:
		return KeyCodeLeft
	case keyboard.KeyArrowRight:
		return KeyCodeRight
	case keyboard.KeyEnter, keyboard.KeySpace:
		return KeyCodeConfirm
	}
	return KeyCodeNone
}

func NewKeypressClientPayload(key KeyCode) Payload {
	rawData, _ := json.Marshal(key)
	return Payload{
		Type: ClientKeypressPayloadType,
		Data: rawData,
	}
}
