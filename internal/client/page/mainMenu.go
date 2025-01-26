package page

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MainMenuOption string

const (
	CreateRoomOption MainMenuOption = "Create Room"
	JoinRoomOption   MainMenuOption = "Join Room"
)

type MainMenu struct {
	cursor  int
	options []MainMenuOption
	msg     string
}

func NewMainMenu() MainMenu {
	return MainMenu{
		cursor: 0,
		options: []MainMenuOption{
			CreateRoomOption,
			JoinRoomOption,
		},
		msg: "",
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.msg = ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return nil, tea.Quit
		case "up", "w":
			m.cursor = m.MoveCursorUp()
		case "down", "s":
			m.cursor = m.MoveCursorDown()
		case "enter", " ":
			return m.getNextPage(), nil
		}
	case error:
		m.msg = "Error: " + msg.Error()
		return m, nil
	}
	return m, nil
}

func (m MainMenu) MoveCursorDown() int {
	return (m.cursor + 1) % len(m.options)
}

func (m MainMenu) MoveCursorUp() int {
	return (m.cursor - 1 + len(m.options)) % len(m.options)
}

func (m MainMenu) getNextPage() tea.Model {
	switch m.options[m.cursor] {
	case CreateRoomOption:
		return NewCreateRoomForm()
	case JoinRoomOption:
		return NewJoinRoomForm()
	}
	return m
}

func (m MainMenu) View() string {
	s := ""
	for i, opt := range m.options {
		if i == m.cursor {
			s += "> "
		} else {
			s += "  "
		}
		s += string(opt) + "\n"
	}
	s += m.msg + "\n"
	return s
}
