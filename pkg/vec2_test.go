package pkg

import "testing"

func TestNewVec2(t *testing.T) {
	x, y := 89, 420
	got := NewVec2(89, 420)
	if got.X != x || got.Y != y {
		t.Errorf("want (%d, %d), got (%d, %d)", x, y, got.X, got.Y)
	}
}

func TestZeroVec2(t *testing.T) {
	got := ZeroVec2()
	if got.X != 0 || got.Y != 0 {
		t.Errorf("want (%d, %d), got (%d, %d)",
			0, 0, got.X, got.Y,
		)
	}
}

func TestVec2_Add(t *testing.T) {
	x1, y1, x2, y2 := 4, -6, 3, 19
	wantX, wantY := x1+x2, y1+y2
	v1 := NewVec2(x1, y1)
	v2 := NewVec2(x2, y2)
	got := v1.Add(v2)
	// check against side effects
	if v1.X != x1 || v1.Y != y1 {
		t.Errorf("v1 should stay at (%d, %d) after Add, got (%d, %d)",
			x1, y1, v1.X, v1.Y,
		)
	}
	if v2.X != x2 || v2.Y != y2 {
		t.Errorf("v2 should stay at (%d, %d) after Add, got (%d, %d)",
			x2, y2, v2.X, v2.Y,
		)
	}
	if got.X != wantX || got.Y != wantY {
		t.Errorf("want (%d, %d), got (%d, %d)",
			wantX, wantY, got.X, got.Y,
		)
	}
}

func TestVec2_Sub(t *testing.T) {
	x1, y1, x2, y2 := 4, -6, 3, 19
	wantX, wantY := x1-x2, y1-y2
	v1 := NewVec2(x1, y1)
	v2 := NewVec2(x2, y2)
	got := v1.Sub(v2)
	// check against side effects
	if v1.X != x1 || v1.Y != y1 {
		t.Errorf("v1 should stay at (%d, %d) after Sub, got (%d, %d)",
			x1, y1, v1.X, v1.Y,
		)
	}
	if v2.X != x2 || v2.Y != y2 {
		t.Errorf("v2 should stay at (%d, %d) after Sub, got (%d, %d)",
			x2, y2, v2.X, v2.Y,
		)
	}
	if got.X != wantX || got.Y != wantY {
		t.Errorf("want (%d, %d), got (%d, %d)",
			wantX, wantY, got.X, got.Y,
		)
	}
}

func TestVec2_Down(t *testing.T) {
	x, y := 4, -6
	wantX, wantY := x, y+1
	v := NewVec2(x, y)
	got := v.Down()
	if got.X != wantX || got.Y != wantY {
		t.Errorf("want (%d, %d), got (%d, %d)",
			wantX, wantY, got.X, got.Y,
		)
	}
	// check against side effects
	if v.X != x || v.Y != y {
		t.Errorf("v should stay at (%d, %d) after Down, got (%d, %d)",
			x, y, v.X, v.Y,
		)
	}
}

func TestVec2_Left(t *testing.T) {
	x, y := 4, -6
	wantX, wantY := x-1, y
	v := NewVec2(x, y)
	got := v.Left()
	if got.X != wantX || got.Y != wantY {
		t.Errorf("want (%d, %d), got (%d, %d)",
			wantX, wantY, got.X, got.Y,
		)
	}
	// check against side effects
	if v.X != x || v.Y != y {
		t.Errorf("v should stay at (%d, %d) after Left, got (%d, %d)",
			x, y, v.X, v.Y,
		)
	}
}

func TestVec2_Right(t *testing.T) {
	x, y := 4, -6
	wantX, wantY := x+1, y
	v := NewVec2(x, y)
	got := v.Right()
	if got.X != wantX || got.Y != wantY {
		t.Errorf("want (%d, %d), got (%d, %d)",
			wantX, wantY, got.X, got.Y,
		)
	}
	// check against side effects
	if v.X != x || v.Y != y {
		t.Errorf("v should stay at (%d, %d) after Right, got (%d, %d)",
			x, y, v.X, v.Y,
		)
	}
}

func TestVec2_Up(t *testing.T) {
	x, y := 4, -6
	wantX, wantY := x, y-1
	v := NewVec2(x, y)
	got := v.Up()
	if got.X != wantX || got.Y != wantY {
		t.Errorf("want (%d, %d), got (%d, %d)",
			wantX, wantY, got.X, got.Y,
		)
	}
	// check against side effects
	if v.X != x || v.Y != y {
		t.Errorf("v should stay at (%d, %d) after Up, got (%d, %d)",
			x, y, v.X, v.Y,
		)
	}
}
