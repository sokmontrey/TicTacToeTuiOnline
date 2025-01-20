package main

func main() {
	//numPlayers := 2
	//url := fmt.Sprintf("http://localhost:4321/create-room?num-players=%d", numPlayers)
	//res, err := http.Get(url)
	//if err != nil {
	//	panic(err)
	//}
	//defer res.Body.Close()
	//log.Println(res.Body)

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
