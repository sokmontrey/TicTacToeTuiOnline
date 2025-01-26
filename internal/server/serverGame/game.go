package serverGame

import "github.com/sokmontrey/TicTacToeTuiOnline/pkg"

type Cell struct {
	position pkg.Vec2
	playerId PlayerId
}

type Game struct {
	players           map[PlayerId]player
	currentPlayerTurn PlayerId
	board             Board
}
