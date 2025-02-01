package pkg

type Vec2 struct {
	X int
	Y int
}

func NewVec2(x int, y int) Vec2 {
	return Vec2{x, y}
}

func ZeroVec2() Vec2 {
	return Vec2{0, 0}
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

func (v Vec2) Up() Vec2 {
	return Vec2{v.X, v.Y - 1}
}

func (v Vec2) Down() Vec2 {
	return Vec2{v.X, v.Y + 1}
}

func (v Vec2) Left() Vec2 {
	return Vec2{v.X - 1, v.Y}
}

func (v Vec2) Right() Vec2 {
	return Vec2{v.X + 1, v.Y}
}

func (v Vec2) UpBy(unit int) Vec2 {
	return Vec2{v.X, v.Y - unit}
}

func (v Vec2) DownBy(unit int) Vec2 {
	return Vec2{v.X, v.Y + unit}
}

func (v Vec2) LeftBy(unit int) Vec2 {
	return Vec2{v.X - unit, v.Y}
}

func (v Vec2) RightBy(unit int) Vec2 {
	return Vec2{v.X + unit, v.Y}
}
