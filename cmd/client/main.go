package main

import "github.com/sokmontrey/TicTacToeTuiOnline/internal/client/page"

func main() {
	var currentPage page.Page = page.NewMainMenu()
	for currentPage != nil {
		currentPage = currentPage.Run()
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
