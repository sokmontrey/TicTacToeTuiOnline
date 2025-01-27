package page

import (
	"encoding/json"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type GameRoom struct {
	pageManager *PageManager
	roomId      string
	displayMsg  string
	conn        *websocket.Conn
	move        chan pkg.Payload
}

func NewGameRoom(pm *PageManager, roomId string) *GameRoom {
	return &GameRoom{
		pageManager: pm,
		roomId:      roomId,
		displayMsg:  "",
		move:        make(chan pkg.Payload),
	}
}

func (m *GameRoom) Init() {
	go m.connectAndListenToServer()
	go m.listenForMove()
}

func (m *GameRoom) Update(msg pageMsg.PageMsg) Command {
	m.displayMsg = "                                          "
	switch msg := msg.(type) {
	case pageMsg.OkMsg:
		m.displayMsg = msg.Data.(string)
	case pageMsg.ErrMsg:
		m.displayMsg = "Error: " + msg.Data.(string)
		return QuitCommand // TODO: handle reconnection
	case pageMsg.KeyMsg:
		switch msg.Key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return QuitCommand
		}
		moveCode := msg.ToMoveCode()
		if moveCode != pkg.MoveCodeNone {
			m.move <- pkg.NewPayload(pkg.ClientMovePayload, moveCode)
		}
		return NoneCommand
	}
	return NoneCommand
}

func (m *GameRoom) View() string {
	s := "Game room\n\n"
	s += fmt.Sprintf("Room id: %s\n", m.roomId)
	s += fmt.Sprintf("%s\n", m.displayMsg)
	return s
}

func (m *GameRoom) listenForMove() {
	for {
		select {
		case move := <-m.move:
			err := move.WsSend(m.conn)
			if err != nil {
				m.pageManager.msg <- pageMsg.NewErrMsg("Unable to send message to the server")
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
		m.pageManager.msg <- pageMsg.NewErrMsg("Unable to connect to the server. Try again later")
		// TODO: handle reconnection
		return
	}
	defer conn.Close()
	m.conn = conn
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			m.pageManager.msg <- pageMsg.NewErrMsg("Connection closed")
			return
		}
		var payload pkg.Payload
		err = json.Unmarshal(msg, &payload)
		if err != nil {
			m.pageManager.msg <- pageMsg.NewErrMsg("Unable to parse server response")
		}
		var msgStr string
		json.Unmarshal(payload.Data, &msgStr)
		if payload.Type == pkg.ServerErrPayload {
			m.pageManager.msg <- pageMsg.NewErrMsg(msgStr)
		} else if payload.Type == pkg.ServerOkPayload {
			m.pageManager.msg <- pageMsg.NewOkMsg(msgStr)
		} // TODO: Game update payload
	}
}
