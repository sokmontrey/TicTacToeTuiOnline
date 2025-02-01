package lobby

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/serverGame"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
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
		game:       serverGame.NewGame(numPlayers, 5),
	}
}

func (r *Room) Start() {
	go r.listenForClientsMove()
}

func (r *Room) listenForClientsMove() {
	for {
		for move := range r.move {
			r.mu.Lock()
			if move.moveCode == payload.MoveCodeNone {
				continue
			}
			var globalPayload, directPayload payload.RawPayload
			if move.moveCode == payload.MoveCodeConfirm {
				if r.IsFull() {
					globalPayload, directPayload = r.game.ConfirmPlayer(move.clientId)
				} else {
					globalPayload = payload.NewNonePayload()
					directPayload = payload.NewErrPayload("Room is not full")
				}
			} else {
				globalPayload, directPayload = r.game.MovePlayer(move.clientId, move.moveCode)
			}
			if globalPayload.Type != payload.NonePayload {
				r.globalBroadcast(globalPayload)
			}
			if directPayload.Type != payload.NonePayload {
				r.directBroadcast(move.clientId, directPayload)
			}
			r.mu.Unlock()
		}
	}
}

func (r *Room) directBroadcast(clientId int, payload payload.RawPayload) {
	r.clients[clientId].SendWs(payload)
}

func (r *Room) globalBroadcast(payload payload.RawPayload) {
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
	r.globalBroadcast(payload.NewJoinedUpdatePayload(clientId))
	r.directBroadcast(clientId,
		payload.NewSyncPayload(
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
	r.globalBroadcast(payload.NewOkPayload(str)) // TODO: left payload
	//TODO: if room is empty, countdown to delete
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
