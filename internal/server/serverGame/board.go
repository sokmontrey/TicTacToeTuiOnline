package serverGame

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type Board struct {
	cells map[pkg.Vec2]PlayerId
}
