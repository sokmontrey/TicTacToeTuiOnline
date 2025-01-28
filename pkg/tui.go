package pkg

import (
	"fmt"
	tc "github.com/gdamore/tcell"
)

type TUI struct {
	CurrentLine int
	Screen      tc.Screen
}

func NewTUI() *TUI {
	screen, _ := tc.NewScreen()
	screen.Init()
	screen.Clear()
	return &TUI{
		CurrentLine: 0,
		Screen:      screen,
	}
}

func (t *TUI) Write(str string) {
	runeStr := []rune(str)
	t.Screen.SetContent(1, t.CurrentLine, runeStr[0], runeStr[1:], tc.StyleDefault)
	t.CurrentLine++
}

func (t *TUI) Writef(format string, args ...any) {
	str := fmt.Sprintf(format, args...)
	runeStr := []rune(str)
	t.Screen.SetContent(1, t.CurrentLine, runeStr[0], runeStr[1:], tc.StyleDefault)
	t.CurrentLine++
}

func (t *TUI) MoveCursor(x int, y int) {
	t.Screen.ShowCursor(x, y)
}

func (t *TUI) SetLine(line int) {
	t.CurrentLine = line
}

func (t *TUI) Up() {
	t.CurrentLine--
}

func (t *TUI) Down() {
	t.CurrentLine++
}

func (t *TUI) Show() {
	t.Screen.Show()
}

func (t *TUI) Clear() {
	t.Screen.Clear()
	t.CurrentLine = 0
}
