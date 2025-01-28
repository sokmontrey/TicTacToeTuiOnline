package pkg

import "github.com/nsf/termbox-go"

func TUIWriteText(line int, str string) {
	strRune := []rune(str)
	for i, ch := range strRune {
		termbox.SetCell(i, line, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
}
