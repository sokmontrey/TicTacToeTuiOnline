package page

import (
	"github.com/eiannone/keyboard"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type PageMsg any

type KeyMsg struct {
	Char rune
	Key  keyboard.Key
}

type ServerPayloadMsg struct {
	Payload pkg.ServerPayload
}

type OkMsg struct {
	Data any
}

type ErrMsg struct {
	Data any
}

func (k KeyMsg) IsChar() bool {
	return k.Char != '\x00'
}

type PageCmd int

const (
	ProgramQuit PageCmd = iota
	NoneCmd
)

type Page interface {
	Init()
	Update(PageMsg) PageCmd
	View() string
}
