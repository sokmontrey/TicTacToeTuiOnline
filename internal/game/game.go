package game

import (
	"github.com/nsf/termbox-go"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type Game struct {
	numPlayers int
	players    map[int]*Player
	radius     int
}

func NewGame(numPlayers int) *Game {
	g := &Game{
		numPlayers: numPlayers,
		players:    make(map[int]*Player),
		radius:     7,
	}
	for i := 1; i <= numPlayers; i++ {
		g.players[i] = NewPlayer(i, pkg.NewVec2(0, 0))
	}
	return g
}

// ============================================================================
// Client Related Functions
// ============================================================================

func (g *Game) UpdatePlayerPosition(playerId int, position pkg.Vec2) {
	g.players[playerId].Position = position
}

func (g *Game) getPlayerCells() map[pkg.Vec2]string {
	playerCells := make(map[pkg.Vec2]string)
	idToMark := []string{"A", "B", "C", "D"}
	for i, player := range g.players {
		playerCells[player.Position] = idToMark[i]
	}
	return playerCells
}

func (g *Game) GetPlayers() map[int]*Player {
	return g.players
}

func (g *Game) GetPlayer(playerId int) *Player {
	return g.players[playerId]
}

func (g *Game) Render(offset int) int {
	playerCells := g.getPlayerCells()
	for y := 0; y < g.radius*2+1; y++ {
		for x := 0; x < g.radius*2+1; x++ {
			cellPos := pkg.NewVec2(x, y)
			mark, ok := playerCells[cellPos.Sub(pkg.NewVec2(g.radius, g.radius))]
			if ok {
				g.setCell(cellPos.DownBy(offset), rune(mark[0]))
			} else {
				g.setCell(cellPos.DownBy(offset), '.')
			}
		}
	}
	return offset + g.radius*2 + 1
}

func (g *Game) setCell(vec2 pkg.Vec2, ch rune) {
	termbox.SetCell(vec2.X*2, vec2.Y, ch, termbox.ColorDefault, termbox.ColorDefault)
}

// ============================================================================
// Server Related Functions
// ============================================================================

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
