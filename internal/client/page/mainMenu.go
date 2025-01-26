package page

import (
	"fmt"
	"github.com/eiannone/keyboard"
)

type MenuOption struct {
	Name   string
	Action func()
}

type MainMenu struct {
	pageManager *PageManager
	options     []MenuOption
	cursor      int
}

func NewMainMenu(pm *PageManager) *MainMenu {
	return &MainMenu{
		pageManager: pm,
		options: []MenuOption{
			{"Create room", pm.ToCreateRoomForm},
			{"Join room", pm.ToJoinRoomForm},
		},
		cursor: 0,
	}
}

func (m *MainMenu) Init() {
}

func (m *MainMenu) Update(msg PageMsg) PageCmd {
	switch msg := msg.(type) {
	case KeyMsg:
		switch msg.Key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return ProgramQuit
		case keyboard.KeyEnter, keyboard.KeySpace:
			m.options[m.cursor].Action()
		case keyboard.KeyArrowUp:
			m.moveCursor(1)
		case keyboard.KeyArrowDown:
			m.moveCursor(-1)
		}

		if !msg.IsChar() {
			return NoneCmd
		}

		switch msg.Char {
		case 'w':
			m.moveCursor(1)
		case 's':
			m.moveCursor(-1)
		}
	}
	return NoneCmd
}

func (m *MainMenu) moveCursor(delta int) {
	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = len(m.options) - 1
	} else if m.cursor >= len(m.options) {
		m.cursor = 0
	}
}

func (m *MainMenu) View() string {
	s := "TicTacToe Online\n\n"
	for i, option := range m.options {
		if i == m.cursor {
			s += fmt.Sprintf("[%s]\n", option.Name)
		} else {
			s += fmt.Sprintf(" %s \n", option.Name)
		}
	}
	return s
}
