package page

import (
	"encoding/json"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type GameRoom struct {
	pageManager *PageManager
	roomId      string
	msg         string
	conn        *websocket.Conn
	moveChan    chan pkg.Payload
}

func NewGameRoom(pm *PageManager, roomId string) *GameRoom {
	return &GameRoom{
		pageManager: pm,
		roomId:      roomId,
		msg:         "",
		moveChan:    make(chan pkg.Payload),
	}
}

func (m *GameRoom) Init() {
	go m.connectAndListenToServer()
	go m.listenForMoves()
}

func (m *GameRoom) Update(msg PageMsg) PageCmd {
	m.msg = ""
	switch msg := msg.(type) {
	case OkMsg:
		m.msg = msg.Data.(string)
	case ErrMsg:
		m.msg = "Error: " + msg.Data.(string)
		return ProgramQuit // TODO: handle reconnection
	case KeyMsg:
		var move pkg.KeyCode
		switch msg.Key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return ProgramQuit
		}
		move = pkg.KeyPressToKeyCode(msg.Key)
		if move != pkg.KeyCodeNone {
			m.moveChan <- pkg.NewKeypressClientPayload(move)
			return NoneCmd
		}
		move = pkg.CharToKeyCode(msg.Char)
		if move != pkg.KeyCodeNone {
			m.moveChan <- pkg.NewKeypressClientPayload(move)
			return NoneCmd
		}
	}
	return NoneCmd
}

func (m *GameRoom) View() string {
	s := "Game room\n\n"
	s += fmt.Sprintf("Room id: %s\n", m.roomId)
	s += fmt.Sprintf("%s\n", m.msg)
	return s
}

func (m *GameRoom) listenForMoves() {
	for {
		select {
		case msg := <-m.moveChan:
			if err := m.conn.WriteJSON(msg); err != nil {
				m.pageManager.msg <- ErrMsg{"Unable to send message to the server"}
				return
			}
		}
	}
}

func (m *GameRoom) connectAndListenToServer() {
	port := "4321"
	url := fmt.Sprintf("ws://localhost:%s/ws/join?room-id=%s", port, m.roomId)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		m.pageManager.msg <- ErrMsg{"Unable to connect to the server. Try again later"}
		// TODO: handle reconnection
		return
	}
	defer conn.Close()
	m.conn = conn
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			m.pageManager.msg <- ErrMsg{"Connection closed"}
			return
		}
		var payload pkg.ServerPayload
		if err := json.Unmarshal(msg, &payload); err != nil {
			m.pageManager.msg <- ErrMsg{"Unable to parse server response"}
		} else if payload.Type == pkg.ServerErrPayloadType {
			m.pageManager.msg <- ErrMsg{payload.Data.(string)}
		} else if payload.Type == pkg.ServerOkPayloadType {
			m.pageManager.msg <- OkMsg{payload.Data.(string)}
		} // TODO: Game update payload
	}
}
