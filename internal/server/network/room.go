package network

import (
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/game"
	"sync"
)

type RoomId string

type GlobalBroadcastState struct {
	Cells        []game.Cell   `json:"cells"`
	PlayerId     game.PlayerId `json:"player-id"`
	IsTerminated bool          `json:"is-terminated"`
}

type Room struct {
	id              RoomId
	numPlayers      int
	game            game.Game
	clients         map[game.PlayerId]*websocket.Conn
	mu              sync.Mutex
	globalBroadcast chan GlobalBroadcastState
}
