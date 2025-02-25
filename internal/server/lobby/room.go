package lobby

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/serverGame"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
	"sync"
	"time"
)

type Room struct {
	id         string
	maxPlayers int
	clients    map[int]*Client
	mu         sync.Mutex
	move       chan ClientMove
	game       *serverGame.Game
	onDelete   func()
}

func NewRoom(numPlayers int, id string, onDelete func()) *Room {
	return &Room{
		id:         id,
		maxPlayers: numPlayers,
		clients:    make(map[int]*Client),
		move:       make(chan ClientMove, 10),
		game:       serverGame.NewGame(numPlayers, 5),
		onDelete:   onDelete,
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
					str := fmt.Sprintf("Room is not full (%d/%d)", r.GetNumClients(), r.maxPlayers)
					directPayload = payload.NewErrPayload(str)
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

func (r *Room) HandleNewConnection(conn *websocket.Conn) payload.RawPayload {
	r.mu.Lock()
	defer r.mu.Unlock()
	clientId := r.AddClient(conn)
	r.globalBroadcast(payload.NewJoinedPayload(clientId))

	r.directBroadcast(clientId,
		payload.NewSyncPayload(
			r.game.GetAllPlayers(),
			r.game.GetAllCells(),
			r.game.GetCurrentTurn(),
			clientId,
		),
	)
	return payload.NewNonePayload()
}

func (r *Room) AddClient(conn *websocket.Conn) int {
	clientId := len(r.clients) + 1
	client := NewClient(clientId, conn, r)
	r.clients[clientId] = client
	client.Run()
	return clientId
}

func (r *Room) RemoveClient(clientId int) {
	r.mu.Lock()
	delete(r.clients, clientId)
	r.mu.Unlock()
	str := fmt.Sprintf("Player %d left the room", clientId)
	r.globalBroadcast(payload.NewErrPayload(str)) // TODO: left payload
	go r.countdownToDelete()
}

func (r *Room) countdownToDelete() {
	time.Sleep(time.Second * 60) // 1 minute to delete
	r.mu.Lock()
	if len(r.clients) == 0 {
		r.onDelete()
	}
	r.mu.Unlock()
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
