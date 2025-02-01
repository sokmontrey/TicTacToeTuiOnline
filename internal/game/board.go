package game

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
)

var checkPairDirs = [4][2]pkg.Vec2{
	{pkg.NewVec2(-1, -1), pkg.NewVec2(1, 1)},
	{pkg.NewVec2(0, -1), pkg.NewVec2(0, 1)},
	{pkg.NewVec2(1, -1), pkg.NewVec2(-1, 1)},
	{pkg.NewVec2(-1, 0), pkg.NewVec2(1, 0)},
}

type Board struct {
	cells map[pkg.Vec2]int
}

func NewBoard() *Board {
	return &Board{
		cells: make(map[pkg.Vec2]int),
	}
}

func (b *Board) GetCell(position pkg.Vec2) int {
	cell, ok := b.cells[position]
	if !ok {
		return -1
	}
	return cell
}

func (b *Board) SetCell(position pkg.Vec2, value int) {
	b.cells[position] = value
}

func (b *Board) IsAdjacent(position pkg.Vec2) bool {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			newPos := position.Add(pkg.NewVec2(dx, dy))
			if b.GetCell(newPos) != -1 {
				return true
			}
		}
	}
	return false
}

func (b *Board) IsEmpty() bool {
	return len(b.cells) == 0
}

func (b *Board) GetAllCells() map[pkg.Vec2]int {
	return b.cells
}

// CheckConnected returns true if the board is connected in any direction with the given number of cells
// Depth-first search from pos outwards in 8 directions
// until we find a cell with a value of num
func (b *Board) CheckConnected(pos pkg.Vec2, numConnect int) map[pkg.Vec2]struct{} {
	cellId := b.GetCell(pos)
	if cellId == -1 {
		return make(map[pkg.Vec2]struct{})
	}

	for _, dirPair := range checkPairDirs {
		countedCells := make(map[pkg.Vec2]struct{})
		countedCells[pos] = struct{}{}

		for _, dir := range dirPair { // only 2 runs
			tempPos := pos.Add(dir)
			for cellId == b.GetCell(tempPos) {
				countedCells[tempPos] = struct{}{}
				tempPos = tempPos.Add(dir)
			}
		}

		if len(countedCells) >= numConnect {
			return countedCells
		}
	}

	return make(map[pkg.Vec2]struct{})
}

func (b *Board) _ConnectedOneDir(pos, dir pkg.Vec2, cellId int) []pkg.Vec2 {
	currCellId := b.GetCell(pos)
	if currCellId == -1 || currCellId != cellId {
		return make([]pkg.Vec2, 0)
	}
	return append(b._ConnectedOneDir(pos.Add(dir), dir, cellId), pos)
}
