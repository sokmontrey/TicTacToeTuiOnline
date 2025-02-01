package serverGame

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/game"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type Game struct {
	numPlayers        int
	players           map[int]*game.Player
	currentTurn       int
	board             *game.Board
	numConnectedCells int
	isRunning         bool
}

func NewGame(numPlayers int, numConnectedCells int) *Game {
	g := &Game{
		numPlayers:        numPlayers,
		players:           make(map[int]*game.Player),
		currentTurn:       1,
		board:             game.NewBoard(),
		numConnectedCells: numConnectedCells,
		isRunning:         true,
	}
	for i := 1; i <= numPlayers; i++ {
		g.players[i] = game.NewPlayer(i, pkg.NewVec2(0, 0))
	}
	return g
}

func (g *Game) MovePlayer(playerId int, moveCode payload.MoveCode) (global payload.RawPayload, direct payload.RawPayload) {
	player := g.players[playerId]
	moveFunc, ok := map[payload.MoveCode]func() pkg.Vec2{
		payload.MoveCodeUp:    player.MoveUp,
		payload.MoveCodeDown:  player.MoveDown,
		payload.MoveCodeLeft:  player.MoveLeft,
		payload.MoveCodeRight: player.MoveRight,
	}[moveCode]
	if ok {
		newPos := moveFunc()
		return payload.NewPositionUpdatePayload(playerId, newPos), payload.NewNonePayload()
	}
	return payload.NewNonePayload(), payload.NewNonePayload()
}

func (g *Game) ConfirmPlayer(playerId int) (global payload.RawPayload, direct payload.RawPayload) {
	if !g.isRunning {
		return payload.NewNonePayload(), payload.NewNonePayload()
	}
	player := g.players[playerId]
	if g.currentTurn != playerId {
		return payload.NewNonePayload(), payload.NewErrPayload("Not your turn!")
	}
	if g.board.GetCell(player.Position) != -1 {
		return payload.NewNonePayload(), payload.NewErrPayload("Cell is already taken!")
	}
	if g.board.IsEmpty() && player.Position != pkg.NewVec2(0, 0) {
		return payload.NewNonePayload(), payload.NewErrPayload("Start at the center!")
	} else if !g.board.IsEmpty() && !g.board.IsAdjacent(player.Position) {
		return payload.NewNonePayload(), payload.NewErrPayload("Too part apart!")
	}
	g.board.SetCell(player.Position, playerId)
	connectedCells := g.board.CheckConnected(player.Position, g.numConnectedCells)
	if len(connectedCells) >= g.numConnectedCells {
		g.isRunning = false
		return payload.NewTerminationPayload(playerId, connectedCells), payload.NewNonePayload()
	}
	g.updateTurn()
	return payload.NewBoardUpdatePayload(player.Position, playerId, g.currentTurn), payload.NewNonePayload()
}

func (g *Game) updateTurn() {
	g.currentTurn++
	if g.currentTurn > g.numPlayers {
		g.currentTurn = 1
	}
}

func (g *Game) GetAllCells() []payload.CellUpdate {
	cellUpdates := make([]payload.CellUpdate, 0)
	for pos, cellId := range g.board.GetAllCells() {
		cellUpdates = append(cellUpdates, payload.CellUpdate{
			CellId:  cellId,
			CellPos: pos,
		})
	}
	return cellUpdates
}

func (g *Game) GetAllPlayers() []payload.PlayerUpdate {
	playerUpdates := make([]payload.PlayerUpdate, 0)
	for id, player := range g.players {
		playerUpdates = append(playerUpdates, payload.PlayerUpdate{
			PlayerId: id,
			Position: player.Position,
		})
	}
	return playerUpdates
}

func (g *Game) GetCurrentTurn() int {
	return g.currentTurn
}
