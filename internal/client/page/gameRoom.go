package page

import (
	"encoding/json"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/clientGame"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type GameRoom struct {
	pageManager *PageManager
	playerId    int
	roomId      string
	displayMsg  string
	conn        *websocket.Conn
	move        chan pkg.Payload
	game        *clientGame.Game
}

func NewGameRoom(pm *PageManager, roomId string) *GameRoom {
	return &GameRoom{
		pageManager: pm,
		playerId:    1,
		roomId:      roomId,
		displayMsg:  "",
		move:        make(chan pkg.Payload),
		game:        clientGame.NewGame(2),
	}
}

func (m *GameRoom) Init() {
	go m.connectAndListenToServer()
	go m.listenForMove()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (m *GameRoom) Render() {
	pkg.TUIWriteText(0, "TicTacToeTui")
	pkg.TUIWriteText(1, fmt.Sprintf("Room id: %s", m.roomId))
	nextLine := m.game.Render(2)
	pkg.TUIWriteText(nextLine, m.displayMsg)
	termbox.Flush()
}

func (m *GameRoom) Update(msg pageMsg.PageMsg) Command {
	m.displayMsg = ""
	switch msg := msg.(type) {
	case pageMsg.JoinedIdMsg:
		m.playerId = msg.PlayerId
		return NoneCommand
	case pageMsg.OkMsg:
		m.displayMsg = msg.Data.(string)
	case pageMsg.ErrMsg:
		m.displayMsg = "Error: " + msg.Data.(string)
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
	case pageMsg.PositionMsg:
		playerId := msg.PlayerId
		position := msg.Position
		m.game.UpdatePlayerPosition(playerId, position)
		if playerId == m.playerId {
			m.game.UpdateCameraPosition(position)
		}
		return NoneCommand
	}
	return NoneCommand
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
		switch payload.Type {
		case pkg.ServerErrPayload:
			m.pageManager.msg <- pageMsg.NewErrMsg(msgStr)
		case pkg.ServerOkPayload:
			m.pageManager.msg <- pageMsg.NewOkMsg(msgStr)
		case pkg.ServerPositionPayload:
			var positionUpdate pkg.PositionUpdate
			json.Unmarshal(payload.Data, &positionUpdate)
			m.pageManager.msg <- pageMsg.NewPositionMsg(positionUpdate.PlayerId, positionUpdate.Position)
		case pkg.ServerJoinedIdPayload:
			var joinedId int
			json.Unmarshal(payload.Data, &joinedId)
			m.pageManager.msg <- pageMsg.NewJoinedIdMsg(joinedId)
		}
		// TODO: Game update payload
	}
}
