package game

import "github.com/sokmontrey/TicTacToeTuiOnline/pkg"

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
