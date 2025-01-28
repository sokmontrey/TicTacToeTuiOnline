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

func (g *Game) getPlayerCells() map[pkg.Vec2]int {
	playerCells := make(map[pkg.Vec2]int)
	for i, player := range g.players {
		playerCells[player.Position] = i
	}
	return playerCells
}

func (g *Game) getPlayerMark(playerId int) (rune, rune, rune) {
	switch playerId {
	case 1:
		return '[', 'X', ']'
	case 2:
		return '(', 'O', ')'
	case 3:
		return '{', '#', '}'
	}
	return '|', '$', '|'
}

func (g *Game) Render(offset int) int {
	g.clearPrevPlayerPositions(offset)
	g.renderPlayerPositions(offset)
	g.clearPrevCells(offset)
	return offset + g.radius*2 + 1
}

func (g *Game) renderPlayerPositions(offset int) {
	playerCells := g.getPlayerCells()
	for y := 0; y < g.radius*2+1; y++ {
		for x := 0; x < g.radius*2+1; x++ {
			cellPos := pkg.NewVec2(x, y)
			id, ok := playerCells[cellPos.Sub(pkg.NewVec2(g.radius, g.radius))]
			left, _, right := g.getPlayerMark(id)
			cellPos = cellPos.DownBy(offset)
			if ok {
				g.drawCursor(cellPos, left, right)
			} else {
				g.setCell(cellPos, '.')
			}
		}
	}
}

func (g *Game) clearPrevPlayerPositions(offset int) {
	for y := 0; y < g.radius*2+1; y++ {
		for x := 0; x < g.radius*2+1; x++ {
			cellPos := pkg.NewVec2(x, y)
			cellPos = cellPos.DownBy(offset)
			g.drawCursor(cellPos, ' ', ' ')
		}
	}
}

func (g *Game) clearPrevCells(offset int) {
	for y := 0; y < g.radius*2+1; y++ {
		for x := 0; x < g.radius*2+1; x++ {
			cellPos := pkg.NewVec2(x, y)
			cellPos = cellPos.DownBy(offset)
			cellPos.X = cellPos.X * 2
			g.setCell(cellPos, '.')
		}
	}
}

func (g *Game) setCell(pos pkg.Vec2, ch rune) {
	termbox.SetCell(pos.X*2+1, pos.Y, ch, termbox.ColorDefault, termbox.ColorDefault)
}

func (g *Game) drawCursor(pos pkg.Vec2, left, right rune) {
	termbox.SetCell(pos.X*2, pos.Y, left, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(pos.X*2+2, pos.Y, right, termbox.ColorDefault, termbox.ColorDefault)
}

// ============================================================================
// Server Related Functions
// ============================================================================

func (g *Game) UpdateState(playerId int, moveCode pkg.MoveCode) (global pkg.Payload, direct pkg.Payload) {
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
	if moveCode == pkg.MoveCodeConfirm {
		// TODO: confirm move
		return pkg.NewNonePayload(), pkg.NewNonePayload()
	}
	return pkg.NewNonePayload(), pkg.NewNonePayload()
}
