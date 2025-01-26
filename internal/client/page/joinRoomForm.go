package page

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type JoinRoomForm struct {
	roomIdInput textinput.Model
	msg         string
}

func NewJoinRoomForm() JoinRoomForm {
	roomIdInp := textinput.New()
	roomIdInp.Prompt = "Enter room id: "
	roomIdInp.Placeholder = "____"
	roomIdInp.CharLimit = 4
	roomIdInp.Focus()

	return JoinRoomForm{
		roomIdInput: roomIdInp,
		msg:         "",
	}
}

func (m JoinRoomForm) Init() tea.Cmd {
	return nil
}

func (m JoinRoomForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.msg = ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return nil, tea.Quit
		case "enter", " ":
			return m.getNextPage(), nil
		}
	case error:
		m.msg = "Error: " + msg.Error()
		return m, nil
	}
	var cmd tea.Cmd
	m.roomIdInput, cmd = m.roomIdInput.Update(msg)
	return m, cmd
}

func (m JoinRoomForm) getNextPage() tea.Model {
	roomId := m.roomIdInput.Value()
	if roomId == "" {
		m.msg = "Please enter a room id"
		return m
	}
	return NewGamePage(roomId)
}

func (m JoinRoomForm) View() string {
	s := "Join an existed room"
	s += m.roomIdInput.View() + "\n"
	s += m.msg + "\n"
	return s
}
