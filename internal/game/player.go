package game

import "github.com/sokmontrey/TicTacToeTuiOnline/pkg"

type Player struct {
	Id       int
	Position pkg.Vec2
}

type PlayerCells map[pkg.Vec2]string

func NewPlayer(id int, position pkg.Vec2) *Player {
	return &Player{
		Id:       id,
		Position: position,
	}
}

func (p *Player) MoveUp() pkg.Vec2 {
	p.Position = p.Position.Up()
	return p.Position
}

func (p *Player) MoveDown() pkg.Vec2 {
	p.Position = p.Position.Down()
	return p.Position
}

func (p *Player) MoveLeft() pkg.Vec2 {
	p.Position = p.Position.Left()
	return p.Position
}

func (p *Player) MoveRight() pkg.Vec2 {
	p.Position = p.Position.Right()
	return p.Position
}
