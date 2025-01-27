package network

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"log"
	"math/rand"
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

var upgrader = websocket.Upgrader{}

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
	// TODO: password for each room
	// TODO: public rooms list
	l.ginEngine.GET("/create-room", l.handleCreateRoom)
	l.ginEngine.GET("/ws/join", l.handleJoin)

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	log.Println("Starting a lobby server on port", port)
	log.Fatal(l.ginEngine.Run(":" + port))
}

// ============================================================================
// Handlers
// ============================================================================

func (l *Lobby) handleJoin(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	roomId := c.Query("room-id")
	room, err := l.GetRoom(roomId)
	stop := l.resWsError(conn, err)
	if stop {
		return
	}

	// also handle full room
	err = room.AddClient(conn)
	stop = l.resWsError(conn, err)
	if stop {
		return
	}

	log.Printf("A new client joined room %s", roomId)
}

// handleCreateRoom handles the httpRequest to create a new room.
// Endpoint example: http://localhost:8080/create-room?num-players=3
func (l *Lobby) handleCreateRoom(c *gin.Context) {
	numPlayersStr := c.Query("num-players")
	numPlayers, err := l.validateNumPlayers(numPlayersStr)
	if l.resHttpError(c, err) {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	id := l.generateRoomId()
	room := NewRoom(numPlayers, id)
	l.rooms[id] = room
	room.Start()
	log.Printf("Created room %s", id)
	payload := pkg.NewPayload(pkg.ServerOkPayload, id)
	payload.HttpSend(http.StatusOK, c)
}

// ============================================================================
// Helpers
// ============================================================================

func (l *Lobby) resWsError(conn *websocket.Conn, err error) bool {
	if err != nil {
		payload := pkg.NewPayload(pkg.ServerErrPayload, err.Error())
		payload.WsSend(conn)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
		return true
	}
	return false
}

func (l *Lobby) resHttpError(c *gin.Context, err error) bool {
	if err != nil {
		payload := pkg.NewPayload(pkg.ServerErrPayload, err.Error())
		payload.HttpSend(http.StatusBadRequest, c)
		return true
	}
	return false
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

func (l *Lobby) generateRoomId() string {
	// TODO: Use a better algorithm for generating room ids with less collisions
	var id string
	ok := true
	for ok {
		id = ""
		for i := 0; i < NumIdDigits; i++ {
			id += strconv.Itoa(rand.Intn(10))
		}
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

func (l *Lobby) GetRoom(roomId string) (*Room, error) {
	_, ok := l.rooms[roomId]
	if !ok {
		return nil, errors.New("room not found")
	}
	return l.rooms[roomId], nil
}
