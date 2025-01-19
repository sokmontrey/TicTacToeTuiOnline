package network

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/game"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"log"
	"sync"
)

type GlobalBroadcastState struct {
	Cells        []game.Cell   `json:"cells"`
	PlayerId     game.PlayerId `json:"player-id"`
	IsTerminated bool          `json:"is-terminated"`
}

type Room struct {
	id              string
	numPlayers      int
	game            *game.Game
	clients         []*websocket.Conn
	globalBroadcast chan GlobalBroadcastState
	mu              sync.Mutex
}

func NewRoom(numPlayers int, id string) *Room {
	return &Room{
		id:              id,
		numPlayers:      numPlayers,
		game:            nil, // TODO: NewGame(),
		clients:         make([]*websocket.Conn, 0, numPlayers),
		globalBroadcast: make(chan GlobalBroadcastState),
	}
}

func (r *Room) Start() {
}

func (r *Room) AddClient(conn *websocket.Conn) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.IsFull() {
		return errors.New("room is full")
	}

	r.clients = append(r.clients, conn)
	go r.handleClient(conn)

	return nil
}

func (r *Room) handleClient(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		r.mu.Lock()
		pkg.SliceRemoveByValue(r.clients, conn)
		r.mu.Unlock()
	}()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		log.Printf("Received message from client: %s", msg)
	}
}

// ============================================================================
// Getters
// ============================================================================

func (r *Room) IsFull() bool {
	return len(r.clients) >= r.numPlayers
}

func (r *Room) GetNumClients() int {
	return len(r.clients)
}
