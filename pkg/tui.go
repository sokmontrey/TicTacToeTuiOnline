package pkg

import "github.com/nsf/termbox-go"

func TUIWriteText(line int, str string) {
	strRune := []rune(str)
	for i, ch := range strRune {
		termbox.SetCell(i, line, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	w, _ := termbox.Size()
	remaining := w - len(strRune)
	for i := 0; i < remaining; i++ {
		termbox.SetCell(i+len(strRune), line, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
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
	for i, ch := range strRune {
		termbox.SetCell(i, line, ch, color, termbox.ColorDefault)
	}
	w, _ := termbox.Size()
	remaining := w - len(strRune)
	for i := 0; i < remaining; i++ {
		termbox.SetCell(i+len(strRune), line, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
}
