package network

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"log"
	"sync"
)

type Room struct {
	id         string
	maxPlayers int
	clients    map[int]*Client
	mu         sync.Mutex
	move       chan clientMove
}

func NewRoom(numPlayers int, id string) *Room {
	return &Room{
		id:         id,
		maxPlayers: numPlayers,
		clients:    make(map[int]*Client),
		move:       make(chan clientMove, 10),
	}
}

func (r *Room) Start() {
	go r.listenForClientsMove()
}

func (r *Room) listenForClientsMove() {
	for {
		for move := range r.move {
			str := fmt.Sprintf("Player %d moved with %v", move.clientId, move.moveCode)
			payload := pkg.NewPayload(pkg.ServerOkPayload, str)
			r.mu.Lock()
			r.broadcast(payload)
			r.mu.Unlock()
		}
	}
}

func (r *Room) broadcast(payload pkg.Payload) {
	for _, client := range r.clients {
		err := payload.WsSend(client.conn)
		if err != nil {
			log.Printf("Error sending payload to client %d: %v", client.clientId, err)
		}
	}
}

func (r *Room) AddClient(conn *websocket.Conn) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.IsFull() {
		return errors.New("room is full")
	}
	clientId := len(r.clients) + 1
	client := NewClient(clientId, conn, r)
	r.clients[clientId] = client
	client.Run()
	str := fmt.Sprintf("Player %d joined the room", clientId)
	payload := pkg.NewPayload(pkg.ServerOkPayload, str)
	r.broadcast(payload)
	return nil
}

func (r *Room) RemoveClient(clientId int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, clientId)
	str := fmt.Sprintf("Player %d left the room", clientId)
	payload := pkg.NewPayload(pkg.ServerOkPayload, str)
	r.broadcast(payload)
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
