package game

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type board struct {
	cells map[pkg.Vec2]PlayerId
}
