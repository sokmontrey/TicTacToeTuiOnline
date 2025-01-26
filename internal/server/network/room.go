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

type GlobalBroadcastState struct {
	Cells        []serverGame.Cell   `json:"cells"`
	PlayerId     serverGame.PlayerId `json:"player-id"`
	IsTerminated bool                `json:"is-terminated"`
}

type Room struct {
	id              string
	maxPlayers      int
	game            *serverGame.Game
	clients         map[*websocket.Conn]serverGame.PlayerId
	globalBroadcast chan GlobalBroadcastState
	mu              sync.Mutex
}

func NewRoom(numPlayers int, id string) *Room {
	return &Room{
		id:              id,
		maxPlayers:      numPlayers,
		game:            nil, // TODO: NewGame(),
		clients:         make(map[*websocket.Conn]serverGame.PlayerId),
		globalBroadcast: make(chan GlobalBroadcastState),
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
	defer r.mu.Unlock()
	return nil
}

func (r *Room) HandleClient(conn *websocket.Conn) {
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
		if r.RoutePayload(conn, payload) {
			return
		}
	}
}

func (r *Room) RoutePayload(conn *websocket.Conn, payload pkg.Payload) bool {
	switch payload.Type {
	case pkg.ClientKeypressPayloadType:
		var keyCode pkg.KeyCode
		err := json.Unmarshal(payload.Data, &keyCode)
		if err != nil {
			log.Printf("Error unmarshaling key code: \"%s\", for room %s", err.Error(), r.id)
			return true
		}
		// TODO: Handle key code
		//log.Printf("Received payload from client: %v", keyCode == pkg.KeyCodeConfirm)
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
