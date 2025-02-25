package pkg

import "math"

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

func (v Vec2) Magnitude() int {
	return int(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v Vec2) Normalize() Vec2 {
	mag := v.Magnitude()
	return Vec2{
		int(math.Round(float64(v.X) / float64(mag))),
		int(math.Round(float64(v.Y) / float64(mag))),
	}
}
