package main

import (
	"encoding/json"
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"log"
)

func test(msgChan chan string) {
	roomId := "6470"
	port := "4321"
	url := fmt.Sprintf("ws://localhost:%s/ws/join?room-id=%s", port, roomId)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			log.Println("Connection closed:", err)
			return
		}

		var payload pkg.ServerPayload
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Println("Invalid server response")
			return
		}

		log.Printf("Received payload: %v", payload)
		msgStr := payload.Data.(string)
		msgChan <- msgStr
	}
}

func testKeyboard(msgChan chan string) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		msgChan <- fmt.Sprintf("You pressed: rune %q, key %X\r\n", char, key)
		if key == keyboard.KeyEsc {
			break
		}
	}
}

func main() {
	serverResponseChan := make(chan string)
	go test(serverResponseChan)

	keyboardChan := make(chan string)
	go testKeyboard(keyboardChan)

	for {
		select {
		case msg := <-serverResponseChan:
			fmt.Println(msg)
		case msg := <-keyboardChan:
			fmt.Println(msg)
		}
	}
}
