package menu

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"net/http"
)

type CreateRoomMenu struct {
	numPlayers    int
	step          int
	msg           string
	createdRoomId string
}

func NewCreateRoomMenu() tea.Model {
	var m tea.Model = CreateRoomMenu{
		numPlayers:    2,
		step:          0,
		msg:           "",
		createdRoomId: "",
	}
	return m
}

func (m CreateRoomMenu) Init() tea.Cmd {
	return nil
}

func (m CreateRoomMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.msg = ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.step <= 0 {
				return NewMainMenu(), nil
			} else {
				m.step--
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "tab", " ":
			m.step++
		}
	case error:
		m.msg = "Error: " + msg.Error()
		return m, nil
	}

	switch m.step {
	case 0:
		m.numPlayers = m.updateNumPlayers(msg)
	default:
		m.step = 0
		newM := m.createRoom(m.numPlayers)
		if newM.createdRoomId != "" {
			JoinRoom(newM.createdRoomId)
			return newM, tea.Quit
		}
		return newM, nil
	}
	return m, nil
}

func (m CreateRoomMenu) updateNumPlayers(msg tea.Msg) int {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "s":
			if m.numPlayers > 2 {
				return m.numPlayers - 1
			}
		case "up", "w":
			// TODO: get max num players from the server with a request
			if m.numPlayers < 4 {
				return m.numPlayers + 1
			}
		}
	}
	return m.numPlayers
}

func (m CreateRoomMenu) View() string {
	s := ""
	s += "Create a new room\n"
	// Number of players
	if m.step == 0 {
		s += "                     ↑\n"
		s += fmt.Sprintf("> Number of players: %d\n", m.numPlayers)
		s += "                     ↓\n"
	} else {
		s += fmt.Sprintf("  Number of players: %d\n", m.numPlayers)
	}
	s += m.msg + "\n"
	return s
}

func (m CreateRoomMenu) createRoom(numPlayers int) CreateRoomMenu {
	PORT := "4321"
	url := fmt.Sprintf("http://localhost:%s/create-room?num-players=%d", PORT, numPlayers)
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		m.step = 0
		m.msg = "Error: Unable to connect to the server. Try again later"
		return m
	}
	defer res.Body.Close()

	var payload pkg.ServerPayload
	if res.StatusCode == http.StatusOK {
		err = json.NewDecoder(res.Body).Decode(&payload)
	}

	if err != nil || payload.Type != pkg.ServerOkPayloadType {
		m.msg = "Error: Unable to connect to the server. Try again later"
		return m
	}

	m.createdRoomId = payload.Data.(string)
	m.msg = "Room created with id: " + m.createdRoomId
	return m
}
