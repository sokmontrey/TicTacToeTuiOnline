package page

import (
	"github.com/eiannone/keyboard"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type MenuOption struct {
	Name   string
	Action func()
}

type MainMenu struct {
	pageManager  *PageManager
	options      []MenuOption
	optionCursor int
	tui          *pkg.TUI
}

func NewMainMenu(pm *PageManager) *MainMenu {
	return &MainMenu{
		pageManager:  pm,
		optionCursor: 0,
		tui:          pkg.NewTUI(),
		options: []MenuOption{
			{"Create a room", pm.ToCreateRoomForm},
			{"Join a room", pm.ToJoinRoomForm},
		},
	}
}

func (m *MainMenu) Init() {
}

func (m *MainMenu) Render() {
	//m.tui.Clear()
	//m.tui.Write("TicTacToe Online")
	//for i, option := range m.options {
	//	if i == m.optionCursor {
	//		m.tui.Writef(">%s", option.Name)
	//	} else {
	//		m.tui.Writef(" %s", option.Name)
	//	}
	//}
	//m.tui.Show()
}

func (m *MainMenu) Update(msg pageMsg.PageMsg) Command {
	switch msg := msg.(type) {
	case pageMsg.KeyMsg:
		switch msg.Key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return QuitCommand
		case keyboard.KeyEnter, keyboard.KeySpace:
			m.options[m.optionCursor].Action()
		case keyboard.KeyArrowUp:
			m.moveCursor(1)
		case keyboard.KeyArrowDown:
			m.moveCursor(-1)
		}

		switch msg.Char {
		case 'w':
			m.moveCursor(1)
		case 's':
			m.moveCursor(-1)
		}
	}
	return NoneCommand
}

func (m *MainMenu) moveCursor(delta int) {
	m.optionCursor += delta
	if m.optionCursor < 0 {
		m.optionCursor = len(m.options) - 1
	} else if m.optionCursor >= len(m.options) {
		m.optionCursor = 0
	}
}
