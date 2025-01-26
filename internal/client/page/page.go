package page

import "github.com/eiannone/keyboard"

type PageMsg any

type KeyMsg struct {
	Char rune
	Key  keyboard.Key
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
