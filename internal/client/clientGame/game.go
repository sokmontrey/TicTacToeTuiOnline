package clientGame

import (
	"github.com/nsf/termbox-go"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/game"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

type PlayerBracket struct {
	Left  rune
	Right rune
}

type PlayerMark rune

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
		board:       game.NewBoard(),
	}
	g.players[1] = game.NewPlayer(1, pkg.NewVec2(1, 0))
	return g
}

func (g *Game) AddNewPlayer() {
	playerId := len(g.players) + 1
	g.players[playerId] = game.NewPlayer(playerId, pkg.NewVec2(playerId, 0))
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

func (g *Game) GetPlayerCursor(playerId int) PlayerBracket {
	var mapIdToMark = map[int]PlayerBracket{
		1: {'[', ']'},
		2: {'(', ')'},
		3: {'{', '}'},
	}
	if mark, ok := mapIdToMark[playerId]; ok {
		return mark
	}
	return PlayerBracket{'|', '|'}
}

func (g *Game) GetPlayerMark(playerId int) PlayerMark {
	var mapIdToMark = map[int]PlayerMark{
		1: 'X',
		2: 'O',
		3: '$',
	}
	if mark, ok := mapIdToMark[playerId]; ok {
		return mark
	}
	return '#'
}

func (g *Game) Render(lineOffset int) int {
	g.rasterScan(lineOffset, func(tuiPos, cellPos pkg.Vec2) {
		if cellPos == pkg.NewVec2(0, 0) {
			g.clearCell(tuiPos, '+')
		} else {
			g.clearCell(tuiPos, '.')
		}
	})
	playerCells := g.getPlayerCells()
	g.rasterScan(lineOffset, func(tuiPos, cellPos pkg.Vec2) {
		id, ok := playerCells[cellPos]
		if ok {
			mark := g.GetPlayerCursor(id)
			g.drawCursor(tuiPos, mark.Left, mark.Right)
		}
	})
	g.rasterScan(lineOffset, func(tuiPos, cellPos pkg.Vec2) {
		cellId := g.board.GetCell(cellPos)
		if cellId != -1 {
			mark := g.GetPlayerMark(cellId)
			g.drawCell(tuiPos, rune(mark))
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

func (g *Game) clearCell(pos pkg.Vec2, mark rune) {
	pos.X = pos.X*2 + 1
	termbox.SetCell(pos.X-1, pos.Y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(pos.X, pos.Y, mark, termbox.ColorDarkGray, termbox.ColorDefault)
	termbox.SetCell(pos.X+1, pos.Y, ' ', termbox.ColorDefault, termbox.ColorDefault)
}

func (g *Game) drawCursor(pos pkg.Vec2, left, right rune) {
	pos.X = pos.X*2 + 1
	termbox.SetCell(pos.X-1, pos.Y, left, termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(pos.X+1, pos.Y, right, termbox.ColorDefault, termbox.ColorDefault)
}

func (g *Game) drawCell(pos pkg.Vec2, mark rune) {
	pos.X = pos.X*2 + 1
	termbox.SetCell(pos.X, pos.Y, mark, termbox.ColorDefault, termbox.ColorDefault)
}

func (g *Game) GetCurrentTurn() int {
	return g.currentTurn
}
