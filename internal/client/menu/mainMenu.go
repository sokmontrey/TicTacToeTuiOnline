package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MainMenu struct {
	cursor  int
	options []string
	actions []func() tea.Model
}

func NewMainMenu() tea.Model {
	var m tea.Model = MainMenu{
		cursor: 0,
		options: []string{
			"Create Room",
			"Join Room",
		},
		actions: []func() tea.Model{
			NewCreateRoomMenu,
			NewJoinRoomMenu,
		},
	}
	return m
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m.actions[m.cursor](), nil // TODO swap Cmd
		}
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
		s += opt + "\n"
	}
	return s
}
