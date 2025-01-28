package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Clear the screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw a simple box
	_, height := termbox.Size()
	startX, startY := 5, 5

	// Add some text inside the box
	ch := 'x'
	termbox.SetCell(startX+2, startY+height/2, ch, termbox.ColorYellow, termbox.ColorDefault)

	// Flush the buffer to display everything
	termbox.Flush()

	// Keep the window open for a few seconds
	time.Sleep(5 * time.Second)
}
