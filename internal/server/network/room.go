package network

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/serverGame"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"log"
	"sync"
)

type Room struct {
	id         string
	maxPlayers int
	game       *serverGame.Game
	clients    map[*websocket.Conn]serverGame.PlayerId
	mu         sync.Mutex
}

func NewRoom(numPlayers int, id string) *Room {
	return &Room{
		id:         id,
		maxPlayers: numPlayers,
		game:       nil, // TODO: NewGame(),
		clients:    make(map[*websocket.Conn]serverGame.PlayerId),
	}
}

func (r *Room) Start() {
}

func (r *Room) AddClient(conn *websocket.Conn) error {
	r.mu.Lock()
	if r.IsFull() {
		return errors.New("room is full")
	}
	r.clients[conn] = serverGame.PlayerId(len(r.clients) + 1)
	r.mu.Unlock()
	go r.listenToClient(conn)
	pkg.NewPayload(pkg.ServerOkPayload, "joined room "+r.id).WsSend(conn)
	return nil
}

func (r *Room) listenToClient(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		r.mu.Lock()
		delete(r.clients, conn)
		// TODO: Broadcast player left, stall the clientGame
		r.mu.Unlock()
	}()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client: \"%s\", for room %s", err.Error(), r.id)
			return
		}
		var payload pkg.Payload
		err = json.Unmarshal(msg, &payload)
		if err != nil {
			log.Printf("Error unmarshaling payload: \"%s\", for room %s", err.Error(), r.id)
			return
		}
		stop := r.RoutePayload(payload)
		if stop {
			return
		}
	}
}

func (r *Room) RoutePayload(payload pkg.Payload) bool {
	switch payload.Type {
	case pkg.ClientMovePayload:
		var moveCode pkg.MoveCode
		err := json.Unmarshal(payload.Data, &moveCode)
		if err != nil {
			log.Printf("Error unmarshaling move code: \"%s\", for room %s", err.Error(), r.id)
			return true
		}
		// TODO: Handle key code
		log.Printf("Received payload from client: %v", moveCode)
	default:
		log.Printf("Received unknown payload from client: %v", payload)
	}
	return false
}

// ============================================================================
// Getters
// ============================================================================

func (r *Room) IsFull() bool {
	return len(r.clients) >= r.maxPlayers
}

func (r *Room) GetNumClients() int {
	return len(r.clients)
}
