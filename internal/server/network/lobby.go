package network

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type Lobby struct {
	rooms map[RoomId]*Room
	mu    sync.Mutex
}

func (l *Lobby) CreateRoom(c *gin.Context) {
}

func (l *Lobby) JoinRoom(c *gin.Context) {

}
