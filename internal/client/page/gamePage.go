package page

import tea "github.com/charmbracelet/bubbletea"

type GamePage struct {
	roomId string
	msg    string
}

func NewGamePage(roomId string) GamePage {
	return GamePage{
		roomId: roomId,
		msg:    "",
	}
}

func (m GamePage) Init() tea.Cmd {
	return nil
}

func (m GamePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.msg = ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return nil, tea.Quit
		case "enter", "tab", " ":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m GamePage) View() string {
	s := "Game page"
	return s
}
