package network

import (
	"encoding/json"
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
	if r.IsFull() {
		return errors.New("room is full")
	}
	r.clients = append(r.clients, conn)
	defer r.mu.Unlock()
	return nil
}

func (r *Room) HandleClient(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		r.mu.Lock()
		pkg.SliceRemoveByValue(r.clients, conn)
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
		r.RoutePayload(conn, payload)
	}
}

func (r *Room) RoutePayload(conn *websocket.Conn, payload pkg.Payload) {
	switch payload.Type {
	case pkg.ReqKeypressPayloadType:
		var keyCode pkg.KeyCode
		err := json.Unmarshal(payload.Data, &keyCode)
		if err != nil {
			log.Printf("Error unmarshaling key code: \"%s\", for room %s", err.Error(), r.id)
			return
		}
		// TODO: Handle key code
		//log.Printf("Received payload from client: %v", keyCode == pkg.KeyCodeConfirm)
	default:
		log.Printf("Received unknown payload from client: %v", payload)
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
