package menu

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"log"
)

type JoinRoomMenu struct {
	roomIdInput textinput.Model
	msg         string
}

func NewJoinRoomMenu() tea.Model {
	roomIdInp := textinput.New()
	roomIdInp.Prompt = "Enter room id: "
	roomIdInp.Placeholder = "____"
	roomIdInp.CharLimit = 4
	roomIdInp.Focus()

	var m tea.Model = JoinRoomMenu{
		roomIdInput: roomIdInp,
		msg:         "",
	}
	return m
}

func JoinRoom(roomId string) {
	url := fmt.Sprintf("ws://localhost:4321/ws/join?room-id=%s", roomId)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			fmt.Println("Connection closed:", err)
			return
		}
		var payload pkg.ServerPayload
		err = json.Unmarshal(msg, &payload)
		if err != nil {
			fmt.Println("Unable to read data from the server")
		}
		log.Print("Log: ")
		log.Println(payload.Type, payload.Data)
	}
}

func (m JoinRoomMenu) Init() tea.Cmd {
	return nil
}

func (m JoinRoomMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return NewMainMenu(), nil
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			value := m.roomIdInput.Value()
			if value != "" {
				JoinRoom(value)
			} else {
				m.msg = "Please enter a room id"
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.roomIdInput, cmd = m.roomIdInput.Update(msg)
	return m, cmd
}

func (m JoinRoomMenu) View() string {
	s := "Join an existed room"
	s += m.roomIdInput.View() + "\n"
	return s
}
