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
	return b.cells[position]
}

func (b *Board) SetCell(position pkg.Vec2, value int) {
	b.cells[position] = value
}
