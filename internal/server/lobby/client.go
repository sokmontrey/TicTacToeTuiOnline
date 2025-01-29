package lobby

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"log"
)

type ClientMove struct {
	clientId int
	moveCode pkg.MoveCode
}

type Client struct {
	clientId int
	conn     *websocket.Conn
	room     *Room
}

func NewClient(clientId int, conn *websocket.Conn, room *Room) *Client {
	return &Client{
		clientId: clientId,
		conn:     conn,
		room:     room,
	}
}

func (c *Client) Run() {
	go c.listenForPayload()
}

func (c *Client) listenForPayload() {
	defer func() {
		c.conn.Close()
		c.room.RemoveClient(c.clientId) // TODO: handle error
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client: \"%s\", for room %s", err.Error(), c.room.id)
			return
		}
		var payload pkg.Payload
		err = json.Unmarshal(msg, &payload)
		if err != nil {
			log.Printf("Error unmarshaling payload: \"%s\", for room %s", err.Error(), c.room.id)
			return
		}
		stop := c.routePayload(payload)
		if stop {
			return
		}
	}
}

func (c *Client) routePayload(payload pkg.Payload) bool {
	switch payload.Type {
	case pkg.ClientMovePayload:
		var moveCode pkg.MoveCode
		err := json.Unmarshal(payload.Data, &moveCode)
		if err != nil {
			log.Printf("Error unmarshaling move code: \"%s\", for room %s", err.Error(), c.room.id)
			return true
		}
		if moveCode != pkg.MoveCodeNone {
			c.room.move <- ClientMove{c.clientId, moveCode}
		}
	}
	return false
}

func (c *Client) SendWs(payload pkg.Payload) {
	err := payload.WsSend(c.conn)
	if err != nil {
		log.Printf("Error sending payload to client %d: %v", c.clientId, err)
	}
}
