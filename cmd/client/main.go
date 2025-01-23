package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sokmontrey/TicTacToeTuiOnline/internal/client/menu"
	"os"
)

func main() {
	var m = menu.NewMainMenu()
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	//if newModel != nil {
	//	p = tea.NewProgram(newModel)
	//}

	//roomId := "2224"
	//url := fmt.Sprintf("ws://localhost:4321/ws/join?room-id=%s", roomId)
	//conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	//if err != nil {
	//	panic(err)
	//}
	//defer conn.Close()
	//
	//for {
	//	msgType, msg, err := conn.ReadMessage()
	//	if err != nil || msgType == websocket.CloseMessage {
	//		fmt.Println("Connection closed:", err)
	//		return
	//	}
	//	fmt.Printf("Received message from server: %s\n", msg)
	//}

	//roomId := got["id"].(string)
	//url = fmt.Sprintf("ws://localhost:1114/ws/join?room-id=%s", roomId)
	//res := wsRequest(t, url)
	//if res.Type != ResponseTypeSuccess && res.Data != "joined room "+roomId {
	//	t.Errorf("got %v, wanted %v", res, NewSuccessResponse("joined room "+roomId))
	//}

	//if err := keyboard.Open(); err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	_ = keyboard.Close()
	//}()
	//
	//fmt.Println("Press ESC to quit")
	//for {
	//	char, key, err := keyboard.GetKey()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
	//	if key == keyboard.KeyEsc {
	//		break
	//	}
	//}
}
