package clientGame

import (
	"github.com/nsf/termbox-go"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/game"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type Game struct {
	numPlayers  int
	players     map[int]*game.Player
	radius      int
	currentTurn int
	cameraPos   pkg.Vec2
	board       *game.Board
}

func NewGame(numPlayers int) *Game {
	g := &Game{
		numPlayers:  numPlayers,
		players:     make(map[int]*game.Player),
		radius:      7,
		currentTurn: 1,
		cameraPos:   pkg.NewVec2(0, 0),
	}
	for i := 1; i <= numPlayers; i++ {
		g.players[i] = game.NewPlayer(i, pkg.NewVec2(0, 0))
	}
	return g
}

func (g *Game) UpdatePlayerPosition(playerId int, position pkg.Vec2) {
	g.players[playerId].Position = position
}

func (g *Game) UpdateCameraPosition(position pkg.Vec2) {
	g.cameraPos = position
}

func (g *Game) UpdateBoard(cellPos pkg.Vec2, cellId int) {
	g.board.SetCell(cellPos, cellId)
}

func (g *Game) UpdateTurn(nextTurn int) {
	g.currentTurn = nextTurn
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
