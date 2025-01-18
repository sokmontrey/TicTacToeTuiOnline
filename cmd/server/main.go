package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/network"
)

func main() {
	PORT := "8080"

	lobby := network.NewLobby()
	r := gin.Default()
	r.GET("/room/create/:numPlayers", lobby.CreateRoom)
	r.GET("/room/join/:room_id", lobby.JoinRoom)

	err := r.Run(":" + PORT)
	if err != nil {
		fmt.Println("Unable to start a network server on port: ", PORT)
	} else {
		fmt.Println("A server is running on port: ", PORT)
	}
}
