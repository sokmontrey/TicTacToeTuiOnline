package clientGame

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Game struct {
	msg  chan string
	conn *websocket.Conn
	mu   sync.Mutex
}

func (g Game) SetConn(conn *websocket.Conn) {
	g.mu.Lock()
	g.conn = conn
	g.mu.Unlock()
}

func NewGame() *Game {
	return &Game{
		msg:  make(chan string),
		conn: nil,
	}
}
