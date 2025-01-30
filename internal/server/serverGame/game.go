package serverGame

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/game"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type Game struct {
	numPlayers  int
	players     map[int]*game.Player
	currentTurn int
	board       *game.Board
}

func NewGame(numPlayers int) *Game {
	g := &Game{
		numPlayers:  numPlayers,
		players:     make(map[int]*game.Player),
		currentTurn: 1,
		board:       game.NewBoard(),
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
	player := g.players[playerId]
	if g.currentTurn != playerId {
		return pkg.NewNonePayload(), pkg.NewPayload(pkg.ServerErrPayload, "Not your turn!")
	}
	if g.board.GetCell(player.Position) != -1 {
		return pkg.NewNonePayload(), pkg.NewPayload(pkg.ServerErrPayload, "Cell is already taken!")
	}
	g.board.SetCell(player.Position, playerId)
	g.updateTurn()
	return pkg.NewBoardUpdatePayload(player.Position, playerId, g.currentTurn), pkg.NewNonePayload()
}

func (g *Game) updateTurn() {
	g.currentTurn++
	if g.currentTurn > g.numPlayers {
		g.currentTurn = 1
	}
}
