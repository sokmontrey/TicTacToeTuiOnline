package main

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/page"
)

const (
	serverAddr = "tictactoetuionline.onrender.com"
)

func main() {
	pageManager := page.NewPageManager(serverAddr)
	pageManager.Init()
	pageManager.ToMainMenu()
	pageManager.Run()
}
