package main

import "github.com/gdamore/tcell"

func main() {
	screen.Clear()
	screen.SetContent(x, 1, 'A', nil, tcell.StyleDefault)
	screen.Show()
}
