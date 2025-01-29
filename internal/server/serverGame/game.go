package serverGame

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/game"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type Game struct {
	numPlayers int
	players    map[int]*game.Player
}

func NewGame(numPlayers int) *Game {
	g := &Game{
		numPlayers: numPlayers,
		players:    make(map[int]*game.Player),
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
	return pkg.NewNonePayload(), pkg.NewNonePayload()
}
