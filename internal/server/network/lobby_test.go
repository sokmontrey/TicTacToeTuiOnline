package network

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"io"
	"net/http"
	"testing"
)

func httpRequest(t *testing.T, url string, expectedStatusCode int) map[string]any {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("got status code %d, wanted %d", resp.StatusCode, expectedStatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func wsRequest(t *testing.T, url string) pkg.Payload {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			t.Log("Connection closed:", err)
			return pkg.Payload{}
		}
		var res pkg.Payload
		err = json.Unmarshal(msg, &res)
		if err != nil {
			t.Fatal(err)
		}
		if res.Type == pkg.ResErrPayloadType || res.Type == pkg.ResOkPayloadType {
			return res
		}
	}
}

func TestLobby_CreateRoom(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1111")

	numPlayers := 3
	url := fmt.Sprintf("http://localhost:1111/create-room?num-players=%d", numPlayers)
	got := httpRequest(t, url, http.StatusOK)

	if len(got["id"].(string)) != NumIdDigits {
		t.Errorf("got id %s, wanted %d digits", got["id"], NumIdDigits)
	}
}

func TestLobby_CreateRoom_InvalidNumPlayers(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1112")

	numPlayers := 0
	url := fmt.Sprintf("http://localhost:1111/create-room?num-players=%d", numPlayers)
	got := httpRequest(t, url, http.StatusBadRequest)
	_, ok := got["error"]
	if !ok {
		t.Errorf("expected error response, got %v", got)
	}
}

func TestLobby_NumRoomsCreated(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1113")

	numPlayers := 3
	url := fmt.Sprintf("http://localhost:1113/create-room?num-players=%d", numPlayers)
	got1 := httpRequest(t, url, http.StatusOK)

	if lobby.CountRooms() != 1 {
		t.Errorf("got %d rooms, wanted 1", lobby.CountRooms())
	}

	url = fmt.Sprintf("http://localhost:1113/create-room?num-players=%d", numPlayers)
	got2 := httpRequest(t, url, http.StatusOK)

	if got1["id"] == got2["id"] {
		t.Errorf("got same id for two consecutive room creations")
	}

	if lobby.CountRooms() != 2 {
		t.Errorf("got %d rooms, wanted 2", lobby.CountRooms())
	}
}

func TestLobby_JoinRoom(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1114")

	numPlayers := 2
	url := fmt.Sprintf("http://localhost:1114/create-room?num-players=%d", numPlayers)
	got := httpRequest(t, url, http.StatusOK)

	roomId := got["id"].(string)
	url = fmt.Sprintf("ws://localhost:1114/ws/join?room-id=%s", roomId)
	res := wsRequest(t, url)
	if res.Type != pkg.ResOkPayloadType && res.Data != "joined room "+roomId {
		t.Errorf("got %v, wanted %v", res, pkg.NewOkResPayload("joined room "+roomId))
	}
}

func TestLobby_JoinRoomInvalidRoomId(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1115")

	numPlayers := 2
	url := fmt.Sprintf("http://localhost:1114/create-room?num-players=%d", numPlayers)
	got := httpRequest(t, url, http.StatusOK)

	roomId := got["id"].(string) + "INVALID"
	url = fmt.Sprintf("ws://localhost:1114/ws/join?room-id=%s", roomId)
	res := wsRequest(t, url)
	if res.Type != pkg.ResErrPayloadType || res.Data != "room not found" {
		t.Errorf("got %v, wanted %v", res, pkg.NewErrResPayload("room not found"))
	}
}

func TestLobby_JoinRoomMaxPlayers(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1116")

	numPlayers := 2
	url := fmt.Sprintf("http://localhost:1114/create-room?num-players=%d", numPlayers)
	got := httpRequest(t, url, http.StatusOK)

	roomId := got["id"].(string)
	url = fmt.Sprintf("ws://localhost:1114/ws/join?room-id=%s", roomId)
	wsRequest(t, url)
	wsRequest(t, url)
	res := wsRequest(t, url)
	if res.Type != pkg.ResErrPayloadType || res.Data != "room is full" {
		t.Errorf("got %v, wanted %v", res, pkg.NewErrResPayload("room is full"))
	}
}
