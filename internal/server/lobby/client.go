package lobby

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/payload"
	"log"
)

type ClientMove struct {
	clientId int
	moveCode payload.MoveCode
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
		var payload payload.RawPayload
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

func (c *Client) routePayload(rawPayload payload.RawPayload) bool {
	switch rawPayload.Type {
	case payload.ClientMoveCodePayload:
		var moveCode payload.MoveCode
		err := json.Unmarshal(rawPayload.Data, &moveCode)
		if err != nil {
			log.Printf("Error unmarshaling move code: \"%s\", for room %s", err.Error(), c.room.id)
			return true
		}
		if moveCode != payload.MoveCodeNone {
			c.room.move <- ClientMove{c.clientId, moveCode}
		}
	}
	return false
}

func (c *Client) SendWs(payload payload.RawPayload) {
	err := payload.WsSend(c.conn)
	if err != nil {
		log.Printf("Error sending payload to client %d: %v", c.clientId, err)
	}
}
