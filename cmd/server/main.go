package main

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/lobby"
)

func main() {
	lobby.NewLobby().Start("4321")
}
