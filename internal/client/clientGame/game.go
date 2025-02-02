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
	numPlayers     int
	players        map[int]*game.Player
	radius         int
	padding        int
	currentTurn    int
	cameraPos      pkg.Vec2
	board          *game.Board
	connectedCells map[pkg.Vec2]struct{}
}

func NewGame(numPlayers int) *Game {
	g := &Game{
		numPlayers:     numPlayers,
		players:        make(map[int]*game.Player),
		radius:         10,
		padding:        1,
		currentTurn:    1,
		cameraPos:      pkg.NewVec2(0, 0),
		board:          game.NewBoard(),
		connectedCells: nil,
	}
	return g
}

func (g *Game) UpdatePlayerPosition(playerId int, pos pkg.Vec2, ownerPlayerId int) {
	player, ok := g.players[playerId]
	if !ok {
		g.players[playerId] = game.NewPlayer(playerId, pos)
		return
	}
	player.Position = pos
	if playerId == ownerPlayerId {
		g.UpdateCameraPosition(pos)
	}
}

func (g *Game) UpdateConnectedCells(connectedCells map[pkg.Vec2]struct{}) {
	g.connectedCells = connectedCells
}

func (g *Game) UpdateCameraPosition(position pkg.Vec2) {
	dir := position.Sub(g.cameraPos)
	if dir.Magnitude() > g.radius-g.padding {
		g.cameraPos = g.cameraPos.Add(dir.Normalize())
	}
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

func (g *Game) GetPlayerColor(playerId int) termbox.Attribute {
	var mapIdToColor = map[int]termbox.Attribute{
		1: termbox.ColorRed,
		2: termbox.ColorBlue,
		3: termbox.ColorGreen,
	}
	if color, ok := mapIdToColor[playerId]; ok {
		return color
	}
	return termbox.ColorDefault
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
			g.clearCell(tuiPos, '+', termbox.ColorDarkGray)
		} else {
			g.clearCell(tuiPos, '.', termbox.ColorDarkGray)
		}
	})
	playerCells := g.getPlayerCells()
	g.rasterScan(lineOffset, func(tuiPos, cellPos pkg.Vec2) {
		id, ok := playerCells[cellPos]
		if ok {
			cursor := g.GetPlayerCursor(id)
			color := g.GetPlayerColor(id)
			g.drawCursor(tuiPos, cursor.Left, cursor.Right, color)
		}
	})
	g.rasterScan(lineOffset, func(tuiPos, cellPos pkg.Vec2) {
		if g.connectedCells != nil {
			_, ok := g.connectedCells[cellPos]
			if ok {
				g.drawCell(tuiPos, '*', termbox.ColorYellow)
				return
			}
		}
		cellId := g.board.GetCell(cellPos)
		if cellId != -1 {
			mark := g.GetPlayerMark(cellId)
			color := g.GetPlayerColor(cellId)
			g.drawCell(tuiPos, rune(mark), color)
		}
	})
	return lineOffset + g.radius*2 + 1
}

func (g *Game) rasterScan(lineOffset int, f func(tuiPos, cellPos pkg.Vec2)) {
	for y := 0; y < 2*g.radius+1; y++ {
		for x := 0; x < 2*g.radius+1; x++ {
			cellPos := pkg.NewVec2(x-g.radius, y-g.radius)
			if cellPos.Magnitude() <= g.radius {
				tuiPos := pkg.NewVec2(x, y+lineOffset)
				f(tuiPos, cellPos.Add(g.cameraPos))
			}
		}
	}
}

func (g *Game) clearCell(pos pkg.Vec2, mark rune, color termbox.Attribute) {
	w, _ := termbox.Size()
	pos.X = pos.X*2 + 1 + w/2 - g.radius*2
	termbox.SetCell(pos.X-1, pos.Y, ' ', color, termbox.ColorDefault)
	termbox.SetCell(pos.X, pos.Y, mark, color, termbox.ColorDefault)
	termbox.SetCell(pos.X+1, pos.Y, ' ', color, termbox.ColorDefault)
}

func (g *Game) drawCursor(pos pkg.Vec2, left, right rune, color termbox.Attribute) {
	w, _ := termbox.Size()
	pos.X = pos.X*2 + 1 + w/2 - g.radius*2
	termbox.SetCell(pos.X-1, pos.Y, left, color, termbox.ColorDefault)
	termbox.SetCell(pos.X+1, pos.Y, right, color, termbox.ColorDefault)
}

func (g *Game) drawCell(pos pkg.Vec2, mark rune, color termbox.Attribute) {
	w, _ := termbox.Size()
	pos.X = pos.X*2 + 1 + w/2 - g.radius*2
	termbox.SetCell(pos.X, pos.Y, mark, color, termbox.ColorDefault)
}

func (g *Game) GetCurrentTurn() int {
	return g.currentTurn
}
