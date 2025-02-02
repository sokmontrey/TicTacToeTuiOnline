package main

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/page"
)

func main() {
	pageManager := page.NewPageManager()
	pageManager.Init()
	pageManager.ToMainMenu()
	pageManager.Run()
}
