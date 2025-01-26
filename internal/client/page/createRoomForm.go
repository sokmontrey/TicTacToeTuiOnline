package page

import (
	"encoding/json"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"net/http"
)

type CreateRoomForm struct {
	numPlayers int
	msg        string
}

func NewCreateRoomForm() CreateRoomForm {
	return CreateRoomForm{
		numPlayers: 2,
		msg:        "",
	}
}

func (m CreateRoomForm) Run() Page {
	p := tea.NewProgram(m)
	finalModel, _ := p.Run()
	if finalModel == nil {
		return nil
	}
	form := finalModel.(CreateRoomForm)
	roomId, err := requestCreateRoom(form.numPlayers)
	if err != nil || roomId == "" {
		fmt.Println("Error: ", err)
		return NewCreateRoomForm()
	}
	return NewGamePage(roomId)
}

func (m CreateRoomForm) Init() tea.Cmd {
	return nil
}

func (m CreateRoomForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.msg = ""
	// TODO: BUG: keyboard input doesn't trigger the first time the form is opened
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return nil, tea.Quit
		case "enter", "tab", " ":
			return m, tea.Quit
		case "up", "w", "d", "right":
			if m.numPlayers < 4 {
				m.numPlayers++
			}
		case "down", "s", "a", "left":
			if m.numPlayers > 2 {
				m.numPlayers--
			}
		}
	case error:
		m.msg = "Error: " + msg.Error()
		return m, nil
	}

	return m, nil
}

func (m CreateRoomForm) View() string {
	s := ""
	s += "Create a new room\n"
	s += fmt.Sprintf("> Number of players: < %d >\n", m.numPlayers)
	s += m.msg + "\n"
	return s
}

func requestCreateRoom(numPlayers int) (string, error) {
	PORT := "4321"
	url := fmt.Sprintf("http://localhost:%s/create-room?num-players=%d", PORT, numPlayers)
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return "", errors.New("unable to connect to the server. Try again later")
	}
	defer res.Body.Close()

	var payload pkg.ServerPayload
	if res.StatusCode == http.StatusOK {
		err = json.NewDecoder(res.Body).Decode(&payload)
	}

	if err != nil || payload.Type != pkg.ServerOkPayloadType {
		return "", errors.New("unable to connect to the server. Try again later")
	}

	return payload.Data.(string), nil
}
