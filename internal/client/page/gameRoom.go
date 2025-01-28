package page

import (
	"encoding/json"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gdamore/tcell"
	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/game"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type GameRoom struct {
	pageManager *PageManager
	playerId    int
	roomId      string
	displayMsg  string
	conn        *websocket.Conn
	move        chan pkg.Payload
	game        *game.Game
	radius      int
}

func NewGameRoom(pm *PageManager, roomId string) *GameRoom {
	screen, _ := tcell.NewScreen()
	screen.Init()
	return &GameRoom{
		pageManager: pm,
		playerId:    1,
		roomId:      roomId,
		displayMsg:  "",
		move:        make(chan pkg.Payload),
		game:        game.NewGame(2),
		radius:      7,
	}
}

func (m *GameRoom) Init() {
	go m.connectAndListenToServer()
	go m.listenForMove()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	//for y := 0; y < m.radius*2+1; y++ {
	//	for x := 0; x < m.radius*2+1; x++ {
	//		m.setCell(pkg.NewVec2(x, y), '.')
	//	}
	//}
}

func (m *GameRoom) Render() {
	//s := "Game room\n\n"
	//s += fmt.Sprintf("Room id: %s\n", m.roomId)
	//
	playerCells := m.game.GetPlayerCells()
	for y := 0; y < m.radius*2+1; y++ {
		for x := 0; x < m.radius*2+1; x++ {
			cellPos := pkg.NewVec2(x, y)
			mark, ok := playerCells[cellPos.Sub(pkg.NewVec2(m.radius, m.radius))]
			if ok {
				m.setCell(cellPos, rune(mark[0]))
			} else {
				m.setCell(cellPos, '.')
			}
		}
	}
	termbox.Flush()
}

func (m *GameRoom) setCell(vec2 pkg.Vec2, ch rune) {
	termbox.SetCell(vec2.X*2, vec2.Y, ch, termbox.ColorDefault, termbox.ColorDefault)
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
	case pageMsg.PositionMsg:
		playerId := msg.PlayerId
		position := msg.Position
		m.game.UpdatePlayerPosition(playerId, position)
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
