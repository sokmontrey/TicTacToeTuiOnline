package serverGame

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/game"
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

func (g *Game) MovePlayer(playerId int, moveCode pkg.MoveCode) (global pkg.Payload, direct pkg.Payload) {
	player := g.players[playerId]
	moveFunc, ok := map[pkg.MoveCode]func() pkg.Vec2{
		pkg.MoveCodeUp:    player.MoveUp,
		pkg.MoveCodeDown:  player.MoveDown,
		pkg.MoveCodeLeft:  player.MoveLeft,
		pkg.MoveCodeRight: player.MoveRight,
	}[moveCode]
	if ok {
		newPos := moveFunc()
		return pkg.NewPositionUpdatePayload(playerId, newPos), pkg.NewNonePayload()
	}
	return pkg.NewNonePayload(), pkg.NewNonePayload()
}

func (g *Game) ConfirmPlayer(playerId int) (global pkg.Payload, direct pkg.Payload) {
	if !g.isRunning {
		return pkg.NewNonePayload(), pkg.NewNonePayload()
	}
	player := g.players[playerId]
	if g.currentTurn != playerId {
		return pkg.NewNonePayload(), pkg.NewPayload(pkg.ServerErrPayload, "Not your turn!")
	}
	if g.board.GetCell(player.Position) != -1 {
		return pkg.NewNonePayload(), pkg.NewPayload(pkg.ServerErrPayload, "Cell is already taken!")
	}
	if g.board.IsEmpty() && player.Position != pkg.NewVec2(0, 0) {
		return pkg.NewNonePayload(), pkg.NewPayload(pkg.ServerErrPayload, "Start at the center!")
	} else if !g.board.IsEmpty() && !g.board.IsAdjacent(player.Position) {
		return pkg.NewNonePayload(), pkg.NewPayload(pkg.ServerErrPayload, "Too part apart!")
	}
	g.board.SetCell(player.Position, playerId)
	connectedCells := g.board.CheckConnected(player.Position, g.numConnectedCells)
	if len(connectedCells) >= g.numConnectedCells {
		g.isRunning = false
		return pkg.NewTerminationPayload(playerId, connectedCells), pkg.NewNonePayload()
	}
	g.updateTurn()
	return pkg.NewBoardUpdatePayload(player.Position, playerId, g.currentTurn), pkg.NewNonePayload()
}

func (g *Game) updateTurn() {
	g.currentTurn++
	if g.currentTurn > g.numPlayers {
		g.currentTurn = 1
	}
}

func (g *Game) GetAllCells() []pkg.CellUpdate {
	cellUpdates := make([]pkg.CellUpdate, 0)
	for pos, cellId := range g.board.GetAllCells() {
		cellUpdates = append(cellUpdates, pkg.CellUpdate{
			CellId:  cellId,
			CellPos: pos,
		})
	}
	return cellUpdates
}

func (g *Game) GetAllPlayers() []pkg.PlayerUpdate {
	playerUpdates := make([]pkg.PlayerUpdate, 0)
	for id, player := range g.players {
		playerUpdates = append(playerUpdates, pkg.PlayerUpdate{
			PlayerId: id,
			Position: player.Position,
		})
	}
	return playerUpdates
}

func (g *Game) GetCurrentTurn() int {
	return g.currentTurn
}
