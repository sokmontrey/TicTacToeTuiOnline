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

func (m MainMenu) Run() Page {
	p := tea.NewProgram(m)
	finalModel, _ := p.Run()
	if finalModel == nil {
		return nil
	}
	finalMainMenu := finalModel.(MainMenu)
	switch finalMainMenu.options[finalMainMenu.cursor] {
	case CreateRoomOption:
		return NewCreateRoomForm()
	case JoinRoomOption:
		//return NewJoinRoomMenu()
	}
	return nil
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
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.options) - 1
			}
		case "down", "s":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter", " ":
			return m, tea.Quit
		}
	case error:
		m.msg = "Error: " + msg.Error()
		return m, nil
	}
	return m, nil
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
