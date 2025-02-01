package page

import (
	"encoding/json"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/clientGame"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/pageMsg"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type GameRoom struct {
	pageManager *PageManager
	playerId    int
	roomId      string
	displayMsg  string
	conn        *websocket.Conn
	move        chan payload.RawPayload
	game        *clientGame.Game
}

func NewGameRoom(pm *PageManager, roomId string) *GameRoom {
	return &GameRoom{
		pageManager: pm,
		playerId:    1,
		roomId:      roomId,
		displayMsg:  "",
		move:        make(chan payload.RawPayload),
		game:        clientGame.NewGame(1),
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
	currentTurn := m.game.GetCurrentTurn()
	if currentTurn == m.playerId {
		pkg.TUIWriteText(2, "Your turn!")
	} else {
		currentTurnMark := m.game.GetPlayerMark(currentTurn)
		bracket := m.game.GetPlayerCursor(currentTurn)
		pkg.TUIWriteText(2, fmt.Sprintf("Turn: %c%c%c", bracket.Left, currentTurnMark, bracket.Right))
	}
	pkg.TUIWriteText(3, m.displayMsg)
	m.game.Render(4)
	termbox.Flush()
}

func (m *GameRoom) Update(msg pageMsg.PageMsg) Command {
	m.displayMsg = ""
	switch msg := msg.(type) {
	case payload.RawPayload:
		switch msg.Type {
		case payload.ServerJoinedPayload:
			p := msg.ToJoinedPayload()
			m.displayMsg = fmt.Sprintf("Player %d joined", p.PlayerId)
		case payload.ServerOkPayload:
			p := msg.ToOkPayload()
			m.displayMsg = p.Value
		case payload.ServerErrPayload:
			p := msg.ToErrPayload()
			m.displayMsg = "Error: " + p.Value
		case payload.ServerPlayerPayload:
			p := msg.ToPlayerPayload()
			m.game.UpdatePlayerPosition(p.PlayerId, p.Position, m.playerId)
		case payload.ServerBoardUpdatePayload:
			p := msg.ToBoardUpdatePayload()
			m.game.UpdateBoard(p.Cell.CellPos, p.Cell.CellId)
			m.game.UpdateTurn(p.NextTurn)
		case payload.ServerTerminationPayload:
			p := msg.ToTerminationPayload()
			m.displayMsg = fmt.Sprintf("Player %d wins!", p.WinnerId)
			m.game.UpdateConnectedCells(p.GetConnectedCellsMap())
		case payload.ServerSyncPayload:
			p := msg.ToSyncUpdatePayload()
			m.playerId = p.CurrentPlayerId
			for _, p := range p.PlayerPositions {
				m.game.UpdatePlayerPosition(p.PlayerId, p.Position, m.playerId)
			}
			for _, c := range p.CellPositions {
				m.game.UpdateBoard(c.CellPos, c.CellId)
			}
			m.game.UpdateTurn(p.CurrentTurn)
		default:
			return NoneCommand
		}
	case pageMsg.KeyMsg:
		switch msg.Key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return QuitCommand
		}
		moveCode := payload.KeyMsgToMoveCode(msg)
		if moveCode != payload.MoveCodeNone {
			m.move <- payload.NewMoveCodePayload(moveCode)
		}
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

		var rawPayload payload.RawPayload
		err = json.Unmarshal(msg, &rawPayload)
		if err != nil {
			m.pageManager.msg <- pageMsg.NewErrMsg("Unable to parse server response")
		}

		if rawPayload.Type == payload.NonePayload {
			continue
		}

		m.pageManager.msg <- rawPayload
	}
}
