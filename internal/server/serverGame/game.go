package serverGame

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type Game struct {
	numPlayers int
	players    map[int]*Player
}

func NewGame(numPlayers int) *Game {
	g := &Game{
		numPlayers: numPlayers,
		players:    make(map[int]*Player),
	}
	for i := 1; i <= numPlayers; i++ {
		g.players[i] = NewPlayer(i, pkg.NewVec2(0, 0))
	}
	return g
}

func (g *Game) UpdateState(playerId int, moveCode pkg.MoveCode) (global pkg.Payload, direct pkg.Payload) {
	player := g.players[playerId]
	updatePosition, ok := map[pkg.MoveCode]func() pkg.Vec2{
		pkg.MoveCodeUp:    player.MoveUp,
		pkg.MoveCodeDown:  player.MoveDown,
		pkg.MoveCodeLeft:  player.MoveLeft,
		pkg.MoveCodeRight: player.MoveRight,
	}[moveCode]
	if ok {
		newPos := updatePosition()
		return pkg.NewPositionUpdatePayload(playerId, newPos), pkg.NewNonePayload()
	}
	if moveCode == pkg.MoveCodeConfirm {
		// TODO: confirm move
		return pkg.NewNonePayload(), pkg.NewNonePayload()
	}
	return pkg.NewNonePayload(), pkg.NewNonePayload()
}
