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
	case payload.JoinedPayload:
		m.displayMsg = fmt.Sprintf("Player %d joined", msg.PlayerId)
	case payload.OkPayload:
		m.displayMsg = msg.Value
	case payload.ErrPayload:
		m.displayMsg = "Error: " + msg.Value
	case payload.PlayerPayload:
		playerId := msg.PlayerId
		position := msg.Position
		m.game.UpdatePlayerPosition(playerId, position, m.playerId)
	case payload.BoardUpdatePayload:
		m.game.UpdateBoard(msg.Cell.CellPos, msg.Cell.CellId)
		m.game.UpdateTurn(msg.NextTurn)
	case payload.TerminationPayload:
		m.displayMsg = fmt.Sprintf("Player %d wins!", msg.WinnerId)
		m.game.UpdateConnectedCells(msg.GetConnectedCellsMap())
	case payload.SyncPayload:
		m.playerId = msg.CurrentPlayerId
		for _, p := range msg.PlayerPositions {
			m.game.UpdatePlayerPosition(p.PlayerId, p.Position, m.playerId)
		}
		for _, c := range msg.CellPositions {
			m.game.UpdateBoard(c.CellPos, c.CellId)
		}
		m.game.UpdateTurn(msg.CurrentTurn)
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

		switch rawPayload.Type {
		case payload.ServerJoinedPayload:
			m.pageManager.msg <- rawPayload.ToJoinedPayload()
		case payload.ServerErrPayload:
			m.pageManager.msg <- rawPayload.ToErrPayload()
		case payload.ServerOkPayload:
			m.pageManager.msg <- rawPayload.ToOkPayload()
		case payload.ServerPlayerUpdatePayload:
			m.pageManager.msg <- rawPayload.ToPlayerPayload()
		case payload.ServerBoardUpdatePayload:
			m.pageManager.msg <- rawPayload.ToBoardUpdatePayload()
		case payload.ServerSyncPayload:
			m.pageManager.msg <- rawPayload.ToSyncUpdatePayload()
		case payload.ServerTerminationPayload:
			m.pageManager.msg <- rawPayload.ToTerminationPayload()
		default:
		}
	}
}
