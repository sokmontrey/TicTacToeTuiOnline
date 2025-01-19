package network

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const (
	MinPlayers  = 2
	MaxPlayers  = 3
	NumIdDigits = 4
)

type Lobby struct {
	rooms     map[string]*Room
	mu        sync.Mutex
	ginEngine *gin.Engine
}

func NewLobby() *Lobby {
	return &Lobby{
		rooms:     make(map[string]*Room),
		ginEngine: gin.Default(),
	}
}

func (l *Lobby) Start(port string) {
	l.ginEngine.GET("/create-room", l.handleCreateRoom)

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	log.Println("Listening the lobby server on port", port)
	log.Fatal(l.ginEngine.Run(":" + port))
}

func (l *Lobby) validateNumPlayers(numPlayersStr string) (int, error) {
	numPlayers, err := strconv.Atoi(numPlayersStr)
	if err != nil {
		return 0, errors.New("invalid num-players")
	}
	if numPlayers < MinPlayers || numPlayers > MaxPlayers {
		return 0, errors.New(fmt.Sprintf("num-players must be between %d and %d", MinPlayers, MaxPlayers))
	}
	return numPlayers, nil
}

// handleCreateRoom handles the request to create a new room.
// Endpoint example: http://localhost:8080/create-room?num-players=3
func (l *Lobby) handleCreateRoom(c *gin.Context) {
	numPlayersStr := c.Query("num-players")
	numPlayers, err := l.validateNumPlayers(numPlayersStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	l.mu.Lock()
	id := l.generateRoomId()
	room := NewRoom(numPlayers, id)
	l.rooms[id] = room
	room.Start()
	c.JSON(http.StatusOK, gin.H{"id": id})
	l.mu.Unlock()
}

func (l *Lobby) generateRoomId() string {
	// TODO: Use a better algorithm for generating room ids with less collisions
	id := pkg.GenerateId(NumIdDigits)
	_, ok := l.rooms[id]
	for ok {
		id = pkg.GenerateId(NumIdDigits)
		_, ok = l.rooms[id]
	}
	return id
}

// ============================================================================
// Getters
// ============================================================================

func (l *Lobby) CountRooms() int {
	return len(l.rooms)
}
