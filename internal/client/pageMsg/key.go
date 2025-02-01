package pageMsg

import "github.com/eiannone/keyboard"

type KeyMsg struct {
	Char rune
	Key  keyboard.Key
}

func NewKeyMsg(char rune, key keyboard.Key) KeyMsg {
	return KeyMsg{char, key}
}
