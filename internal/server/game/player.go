package game

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type PlayerId byte

type player struct {
	id       PlayerId
	position pkg.Vec2
}
