package game

import (
	"github.com/nsf/termbox-go"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type Game struct {
	numPlayers int
	players    map[int]*Player
	radius     int
	cameraPos  pkg.Vec2
}

func NewGame(numPlayers int) *Game {
	g := &Game{
		numPlayers: numPlayers,
		players:    make(map[int]*Player),
		radius:     7,
		cameraPos:  pkg.NewVec2(0, 0),
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

func (g *Game) UpdateCameraPosition(position pkg.Vec2) {
	g.cameraPos = position
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

func (g *Game) Render(lineOffset int) int {
	g.rasterScan(lineOffset, func(tuiPos, cellPos pkg.Vec2) {
		g.resetCell(tuiPos)
	})
	playerCells := g.getPlayerCells()
	g.rasterScan(lineOffset, func(tuiPos, cellPos pkg.Vec2) {
		id, ok := playerCells[cellPos]
		if ok {
			left, _, right := g.getPlayerMark(id)
			g.drawCursor(tuiPos, left, right)
		}
	})
	return lineOffset + g.radius*2 + 1
}

func (g *Game) rasterScan(lineOffset int, f func(tuiPos, cellPos pkg.Vec2)) {
	for y := 0; y < 2*g.radius+1; y++ {
		for x := 0; x < 2*g.radius+1; x++ {
			cellPos := pkg.NewVec2(x-g.radius, y-g.radius)
			cellPos = cellPos.Add(g.cameraPos)
			f(pkg.NewVec2(x, y+lineOffset), cellPos)
		}
	}
}

func (g *Game) resetCell(pos pkg.Vec2) {
	pos.X = pos.X*2 + 1
	termbox.SetCell(pos.X-1, pos.Y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(pos.X, pos.Y, '.', termbox.ColorDarkGray, termbox.ColorDefault)
	termbox.SetCell(pos.X+1, pos.Y, ' ', termbox.ColorDefault, termbox.ColorDefault)
}

func (g *Game) drawCursor(pos pkg.Vec2, left, right rune) {
	pos.X = pos.X*2 + 1
	termbox.SetCell(pos.X-1, pos.Y, left, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(pos.X+1, pos.Y, right, termbox.ColorDefault, termbox.ColorDefault)
}

// ============================================================================
// Server Related Functions
// ============================================================================

func (g *Game) ConfirmPlayer(playerId int) (global pkg.Payload, direct pkg.Payload) {
	return pkg.NewNonePayload(), pkg.NewNonePayload()
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
	if moveCode == pkg.MoveCodeConfirm {
		// TODO: confirm move
		return pkg.NewNonePayload(), pkg.NewNonePayload()
	}
	return pkg.NewNonePayload(), pkg.NewNonePayload()
}
