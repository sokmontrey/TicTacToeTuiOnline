package serverGame

import "github.com/sokmontrey/TicTacToeTuiOnline/pkg"

type Player struct {
	id       int
	position pkg.Vec2
}

func NewPlayer(id int, position pkg.Vec2) *Player {
	return &Player{
		id:       id,
		position: position,
	}
}

func (p *Player) MoveUp() pkg.Vec2 {
	p.position = p.position.Up()
	return p.position
}

func (p *Player) MoveDown() pkg.Vec2 {
	p.position = p.position.Down()
	return p.position
}

func (p *Player) MoveLeft() pkg.Vec2 {
	p.position = p.position.Left()
	return p.position
}

func (p *Player) MoveRight() pkg.Vec2 {
	p.position = p.position.Right()
	return p.position
}
