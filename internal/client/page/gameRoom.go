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
		game:        clientGame.NewGame(1),
	}
}

func (m *GameRoom) Init() {
	go m.connectAndListenToServer()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (m *GameRoom) Render() {
	pkg.TUIWriteText(1, "TicTac ToeTui")
	pkg.TUIWriteText(2, fmt.Sprintf("Room id: %s", m.roomId))
	currentTurn := m.game.GetCurrentTurn()
	if currentTurn == m.playerId {
		pkg.TUIWriteTextWithColor(4, "YOUR TURN", termbox.ColorGreen)
	} else {
		currentTurnMark := m.game.GetPlayerMark(currentTurn)
		bracket := m.game.GetPlayerCursor(currentTurn)
		color := m.game.GetPlayerColor(currentTurn)
		symbol := fmt.Sprintf("%c%c%c", bracket.Left, currentTurnMark, bracket.Right)
		pkg.TUIWriteTextWithColor(4, symbol, color)
	}
	nextLine := m.game.Render(5)
	pkg.TUIWriteTextWithColor(nextLine, m.displayMsg, termbox.ColorRed)
	pkg.TUIWriteTextWithColor(nextLine+1, "  W  ", termbox.ColorLightGray)
	pkg.TUIWriteTextWithColor(nextLine+2, "A S D", termbox.ColorLightGray)
	pkg.TUIWriteTextWithColor(nextLine+3, "SPACE", termbox.ColorLightGray)
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
		if moveCode == payload.MoveCodeNone {
			return NoneCommand
		}
		err := payload.NewMoveCodePayload(moveCode).WsSend(m.conn)
		if err != nil {
			m.displayMsg = "Unable to send message to the server"
		}
	}
	return NoneCommand
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
