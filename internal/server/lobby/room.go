package lobby

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/serverGame"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"sync"
)

type Room struct {
	id         string
	maxPlayers int
	clients    map[int]*Client
	mu         sync.Mutex
	move       chan ClientMove
	game       *serverGame.Game
}

func NewRoom(numPlayers int, id string) *Room {
	return &Room{
		id:         id,
		maxPlayers: numPlayers,
		clients:    make(map[int]*Client),
		move:       make(chan ClientMove, 10),
		game:       serverGame.NewGame(numPlayers),
	}
}

func (r *Room) Start() {
	go r.listenForClientsMove()
}

func (r *Room) listenForClientsMove() {
	for {
		for move := range r.move {
			r.mu.Lock()
			if move.moveCode == pkg.MoveCodeNone {
				continue
			}
			var globalPayload, directPayload pkg.Payload
			if move.moveCode == pkg.MoveCodeConfirm {
				if r.IsFull() {
					globalPayload, directPayload = r.game.ConfirmPlayer(move.clientId)
				} else {
					globalPayload = pkg.NewNonePayload()
					directPayload = pkg.NewPayload(pkg.ServerErrPayload, "Room is not full")
				}
			} else {
				globalPayload, directPayload = r.game.MovePlayer(move.clientId, move.moveCode)
			}
			if globalPayload.Type != pkg.NonePayload {
				r.globalBroadcast(globalPayload)
			}
			if directPayload.Type != pkg.NonePayload {
				r.directBroadcast(move.clientId, directPayload)
			}
			r.mu.Unlock()
		}
	}
}

func (r *Room) directBroadcast(clientId int, payload pkg.Payload) {
	r.clients[clientId].SendWs(payload)
}

func (r *Room) globalBroadcast(payload pkg.Payload) {
	for _, client := range r.clients {
		client.SendWs(payload)
	}
}

func (r *Room) AddClient(conn *websocket.Conn) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.IsFull() {
		return errors.New("room is full")
	}
	clientId := r.CreateClient(conn)
	r.globalBroadcast(pkg.NewJoinedUpdatePayload(clientId))
	r.directBroadcast(clientId,
		pkg.NewSyncPayload(
			r.game.GetAllPlayers(),
			r.game.GetAllCells(),
			r.game.GetCurrentTurn(),
			clientId,
		),
	)
	return nil
}

func (r *Room) CreateClient(conn *websocket.Conn) int {
	clientId := len(r.clients) + 1
	client := NewClient(clientId, conn, r)
	r.clients[clientId] = client
	client.Run()
	return clientId
}

func (r *Room) RemoveClient(clientId int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, clientId)
	str := fmt.Sprintf("Player %d left the room", clientId)
	payload := pkg.NewPayload(pkg.ServerOkPayload, str)
	r.globalBroadcast(payload)
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
