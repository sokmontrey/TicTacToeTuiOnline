package network

import (
	"github.com/gin-gonic/gin"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/server/utils"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"strconv"
	"sync"
)

const (
	maxPlayers  = 3
	numIdDigits = 4
)

type Lobby struct {
	rooms map[string]*Room
	mu    sync.Mutex
}

func NewLobby() *Lobby {
	return &Lobby{
		rooms: make(map[string]*Room),
	}
}

func (l *Lobby) GenerateId(digits int) string {
	id := pkg.GenerateId(digits)
	_, ok := l.rooms[id]
	for ok {
		id = pkg.GenerateId(digits)
		_, ok = l.rooms[id]
	}
	return id
}

func (l *Lobby) CreateRoom(c *gin.Context) {
	responder := utils.NewHttpResponder(c)
	numPlayersParam := c.Param("numPlayers")
	numPlayers, err := strconv.Atoi(numPlayersParam)
	if err != nil {
		responder.NumPlayersParamError(numPlayersParam)
	} else if numPlayers > maxPlayers {
		responder.MaxNumPlayerError(maxPlayers)
	}
	l.mu.Lock()
	id := l.GenerateId(numIdDigits)
	l.rooms[id] = NewRoom(numPlayers, id)
	go l.rooms[id].Run()
	responder.RoomCreated(id)
	l.mu.Unlock()
}

func (l *Lobby) JoinRoom(c *gin.Context) {

}
