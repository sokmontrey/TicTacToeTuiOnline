package main

import (
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/network"
)

func main() {
	lobby := network.NewLobby()
	lobby.Start("4321")
}
