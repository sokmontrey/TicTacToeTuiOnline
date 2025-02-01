package page

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
)

type Command int

const (
	QuitCommand Command = iota
	NoneCommand
)

type Page interface {
	Init()
	Update(msg pageMsg.PageMsg) Command
	Render()
}
