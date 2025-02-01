package main

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/lobby"
)

func main() {
	lobby := lobby.NewLobby()
	lobby.Start("4321")
}
