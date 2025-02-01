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
	case payload.JoinedUpdate:
		m.displayMsg = fmt.Sprintf("Player %d joined", msg.PlayerId)
	case pageMsg.SyncMsg:
		m.playerId = msg.CurrentPlayerId
		for _, p := range msg.PlayerPositions {
			m.game.UpdatePlayerPosition(p.PlayerId, p.Position)
			if p.PlayerId == m.playerId {
				m.game.UpdateCameraPosition(p.Position)
			}
		}
		for _, c := range msg.CellPositions {
			m.game.UpdateBoard(c.CellPos, c.CellId)
		}
		m.game.UpdateTurn(msg.CurrentTurn)
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
		if moveCode != payload.MoveCodeNone {
			m.move <- payload.NewPayload(payload.ClientMovePayload, moveCode)
		}
	case pageMsg.PlayerPositionMsg:
		playerId := msg.PlayerId
		position := msg.Position
		m.game.UpdatePlayerPosition(playerId, position)
		if playerId == m.playerId {
			m.game.UpdateCameraPosition(position)
		}
	case pageMsg.BoardUpdateMsg:
		m.game.UpdateBoard(msg.CellPos, msg.CellId)
		m.game.UpdateTurn(msg.NextTurn)
	case pageMsg.TerminationMsg:
		m.displayMsg = fmt.Sprintf("Player %d wins!", msg.WinnerId)
		m.game.UpdateConnectedCells(msg.ConnectedCells)
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
			m.pageManager.msg <- rawPayload.ToJoinedUpdate()
		case payload.ServerErrPayload:
			var msgStr string
			json.Unmarshal(rawPayload.Data, &msgStr)
			m.pageManager.msg <- pageMsg.NewErrMsg(msgStr)
		case payload.ServerOkPayload:
			var msgStr string
			json.Unmarshal(rawPayload.Data, &msgStr)
			m.pageManager.msg <- pageMsg.NewOkMsg(msgStr)
		case payload.ServerPositionPayload:
			var positionUpdate payload.PlayerUpdate
			json.Unmarshal(rawPayload.Data, &positionUpdate)
			m.pageManager.msg <- pageMsg.NewPositionMsg(positionUpdate.PlayerId, positionUpdate.Position)
		case payload.ServerSyncPayload:
			var syncData payload.SyncUpdate
			json.Unmarshal(rawPayload.Data, &syncData)
			m.pageManager.msg <- pageMsg.NewSyncMsg( // TODO; for each of these, pass payload data directly
				syncData.PlayerPositions,
				syncData.CellPositions,
				syncData.CurrentTurn,
				syncData.CurrentPlayerId,
			)
		case payload.ServerBoardUpdatePayload:
			var boardUpdate payload.BoardUpdate
			json.Unmarshal(rawPayload.Data, &boardUpdate)
			m.pageManager.msg <- pageMsg.NewBoardUpdateMsg(
				boardUpdate.Cell.CellPos,
				boardUpdate.Cell.CellId,
				boardUpdate.NextTurn,
			)
		case payload.ServerTerminationPayload:
			var terminationUpdate payload.TerminationUpdate
			json.Unmarshal(rawPayload.Data, &terminationUpdate)
			m.pageManager.msg <- pageMsg.NewTerminationMsg(
				terminationUpdate.WinnerId,
				terminationUpdate.ConnectedCells,
			)
		}
	}
}
