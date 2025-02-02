package page

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/nsf/termbox-go"
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
	msg          string
}

func NewMainMenu(pm *PageManager) *MainMenu {
	return &MainMenu{
		pageManager:  pm,
		optionCursor: 0,
		options: []MenuOption{
			{"Create a room", pm.ToCreateRoomForm},
			{"Join a room", pm.ToJoinRoomForm},
		},
		msg: "",
	}
}

func (m *MainMenu) SetMsg(msg string) {
	m.msg = msg
}

func (m *MainMenu) Init() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (m *MainMenu) Render() {
	pkg.TUIWriteText(1, "TicTacToe Tui")
	for i, option := range m.options {
		if i == m.optionCursor {
			pkg.TUIWriteTextWithColor(2+i, fmt.Sprintf("[%s]", option.Name), termbox.ColorGreen)
		} else {
			pkg.TUIWriteTextWithColor(2+i, fmt.Sprintf(" %s ", option.Name), termbox.ColorDefault)
		}
	}
	pkg.TUIWriteTextWithColor(len(m.options)+2, m.msg, termbox.ColorRed)
	termbox.Flush()
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
