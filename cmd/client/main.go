package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/page"
)

func main() {
	mainMenu := page.NewMainMenu()
	p := tea.NewProgram(page.NewPageManager(mainMenu))
	if _, err := p.Run(); err != nil {
		fmt.Println("There has been an error: ", err)
		fmt.Println("Exiting...")
		return
	}

	//if err := keyboard.Open(); err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	_ = keyboard.Close()
	//}()
	//
	//fmt.Println("Press ESC to quit")
	//for {
	//	char, key, err := keyboard.GetKey()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
	//	if key == keyboard.KeyEsc {
	//		break
	//	}
	//}
}
