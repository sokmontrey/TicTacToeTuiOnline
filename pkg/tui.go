package pkg

import "github.com/nsf/termbox-go"

func TUIWriteText(line int, str string) {
	TUIWriteTextWithColor(line, str, termbox.ColorDefault)
}

func TUILine(line int) {
	w, _ := termbox.Size()
	remaining := w
	for i := 0; i < remaining; i++ {
		termbox.SetCell(i, line, '-', termbox.ColorDarkGray, termbox.ColorDefault)
	}
}

func TUIWriteTextWithColor(line int, str string, color termbox.Attribute) {
	strRune := []rune(str)
	w, _ := termbox.Size()
	remaining := (w - len(strRune)) / 2

	for i := 0; i < remaining; i++ {
		termbox.SetCell(i, line, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}

	for i, ch := range strRune {
		termbox.SetCell(i+remaining, line, ch, color, termbox.ColorDefault)
	}

	for i := 0; i < remaining; i++ {
		termbox.SetCell(i+len(strRune)+remaining, line, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
}
