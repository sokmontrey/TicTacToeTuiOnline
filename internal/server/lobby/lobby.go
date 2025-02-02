package lobby

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
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
	room, rp := l.GetRoom(roomId)
	if l.resWsError(conn, rp) {
		log.Printf("Client tried to join room %s, but it is not found", roomId)
		return
	}

	if room.IsFull() {
		log.Printf("Client tried to join room %s, but it is full", roomId)
		l.resWsError(conn, payload.NewClosePayload("Room is full"))
		return
	}

	rp = room.HandleNewConnection(conn)
	if l.resWsError(conn, rp) {
		return
	}

	log.Printf("A new client joined room %s", roomId)
}

// handleCreateRoom handles the httpRequest to create a new room.
// Endpoint example: http://localhost:8080/create-room?num-players=3
func (l *Lobby) handleCreateRoom(c *gin.Context) {
	numPlayersStr := c.Query("num-players")
	numPlayers, rp := l.validateNumPlayers(numPlayersStr)
	if l.resHttpError(c, rp) {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	id := l.generateRoomId()
	room := NewRoom(numPlayers, id)
	l.rooms[id] = room
	room.Start()
	log.Printf("Created room %s", id)
	payload.NewOkPayload(id).HttpSend(http.StatusOK, c)
}

// ============================================================================
// Helpers
// ============================================================================

func (l *Lobby) resWsError(conn *websocket.Conn, rp payload.RawPayload) bool {
	if rp.Type != payload.NonePayload {
		rp.WsSend(conn)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""))
		conn.Close()
		return true
	}
	return false
}

func (l *Lobby) resHttpError(c *gin.Context, rp payload.RawPayload) bool {
	if rp.Type != payload.NonePayload {
		rp.HttpSend(http.StatusBadRequest, c)
		return true
	}
	return false
}

func (l *Lobby) validateNumPlayers(numPlayersStr string) (int, payload.RawPayload) {
	numPlayers, err := strconv.Atoi(numPlayersStr)
	if err != nil {
		return 0, payload.NewErrPayload("Invalid num-players")
	}
	if numPlayers < MinPlayers || numPlayers > MaxPlayers {
		str := fmt.Sprintf("num-players must be between %d and %d", MinPlayers, MaxPlayers)
		return 0, payload.NewErrPayload(str)
	}
	return numPlayers, payload.NewNonePayload()
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

func (l *Lobby) GetRoom(roomId string) (*Room, payload.RawPayload) {
	_, ok := l.rooms[roomId]
	if !ok {
		return nil, payload.NewClosePayload("Room not found")
	}
	return l.rooms[roomId], payload.NewNonePayload()
}
