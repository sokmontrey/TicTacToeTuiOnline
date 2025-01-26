package main

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/page"
	"log"
)

func main() {
	pageManager := page.NewPageManager()
	pageManager.ToMainMenu()
	log.Println("Starting the program...")
	pageManager.Init()
	pageManager.Run()

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
