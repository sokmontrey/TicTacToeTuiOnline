package page

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
	"net/http"
)

type CreateRoomForm struct {
	maxPlayers  int
	minPlayers  int
	numPlayers  int
	pageManager *PageManager
	msg         string
}

func NewCreateRoomForm(pm *PageManager) *CreateRoomForm {
	return &CreateRoomForm{
		maxPlayers:  4,
		minPlayers:  2,
		numPlayers:  2,
		pageManager: pm,
		msg:         "",
	}
}

func (m *CreateRoomForm) Init() {
}

func (m *CreateRoomForm) Render() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Create room")
	fmt.Printf("Number of players: < %d >\n", m.numPlayers)
	fmt.Printf("%s\n", m.msg)
}

func (m *CreateRoomForm) Update(msg pageMsg.PageMsg) Command {
	m.msg = ""
	switch msg := msg.(type) {
	case pageMsg.KeyMsg:
		switch msg.Key {
		case keyboard.KeyCtrlC:
			return QuitCommand
		case keyboard.KeyEsc:
			m.pageManager.ToMainMenu()
		case keyboard.KeyArrowUp, keyboard.KeyArrowRight:
			m.updateNumPlayers(1)
		case keyboard.KeyArrowDown, keyboard.KeyArrowLeft:
			m.updateNumPlayers(-1)
		case keyboard.KeyEnter:
			roomId, err := m.requestCreateRoom()
			if err != nil {
				m.msg = "Error: " + err.Error()
			} else {
				m.pageManager.ToGameRoom(roomId)
			}
			return NoneCommand
		}

		switch msg.Char {
		case 'w', 'd':
			m.updateNumPlayers(1)
		case 's', 'a':
			m.updateNumPlayers(-1)
		}
	}
	return NoneCommand
}

func (m *CreateRoomForm) updateNumPlayers(delta int) {
	m.numPlayers += delta
	if m.numPlayers < m.minPlayers {
		m.numPlayers = m.minPlayers
	} else if m.numPlayers > m.maxPlayers {
		m.numPlayers = m.maxPlayers
	}
}

func (m *CreateRoomForm) requestCreateRoom() (string, error) {
	PORT := "4321"
	url := fmt.Sprintf("http://localhost:%s/create-room?num-players=%d", PORT, m.numPlayers)
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return "", errors.New("unable to connect to the server. Try again later")
	}
	defer res.Body.Close()
	var rawPayload payload.RawPayload
	if res.StatusCode == http.StatusOK {
		err = json.NewDecoder(res.Body).Decode(&rawPayload)
	}
	if err != nil {
		return "", errors.New("unable to connect to the server. Try again later")
	}
	roomId := rawPayload.ToOkPayload().Value
	return roomId, err
}
