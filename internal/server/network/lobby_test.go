package network

import (
	"fmt"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"net/http"
	"testing"
)

func TestLobby_CreateRoom(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1111")

	numPlayers := 3
	url := fmt.Sprintf("http://localhost:1111/create-room?num-players=%d", numPlayers)
	got := HttpRequest(t, url, http.StatusOK)

	if len(got["id"].(string)) != NumIdDigits {
		t.Errorf("got id %s, wanted %d digits", got["id"], NumIdDigits)
	}
}

func TestLobby_CreateRoom_InvalidNumPlayers(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1112")

	numPlayers := 0
	url := fmt.Sprintf("http://localhost:1111/create-room?num-players=%d", numPlayers)
	got := HttpRequest(t, url, http.StatusBadRequest)
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
	got1 := HttpRequest(t, url, http.StatusOK)

	if lobby.CountRooms() != 1 {
		t.Errorf("got %d rooms, wanted 1", lobby.CountRooms())
	}

	url = fmt.Sprintf("http://localhost:1113/create-room?num-players=%d", numPlayers)
	got2 := HttpRequest(t, url, http.StatusOK)

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
	got := HttpRequest(t, url, http.StatusOK)

	roomId := got["id"].(string)
	url = fmt.Sprintf("ws://localhost:1114/ws/join?room-id=%s", roomId)
	res := WsRequest(t, url)
	if res.Type != pkg.ServerOkPayloadType && res.Data != "joined room "+roomId {
		t.Errorf("got %v, wanted %v", res, pkg.NewOkServerPayload("joined room "+roomId))
	}
}

func TestLobby_JoinRoomInvalidRoomId(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1115")

	numPlayers := 2
	url := fmt.Sprintf("http://localhost:1114/create-room?num-players=%d", numPlayers)
	got := HttpRequest(t, url, http.StatusOK)

	roomId := got["id"].(string) + "INVALID"
	url = fmt.Sprintf("ws://localhost:1114/ws/join?room-id=%s", roomId)
	res := WsRequest(t, url)
	if res.Type != pkg.ServerErrPayloadType || res.Data != "room not found" {
		t.Errorf("got %v, wanted %v", res, pkg.NewErrServerPayload("room not found"))
	}
}

func TestLobby_JoinRoomMaxPlayers(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1116")

	numPlayers := 2
	url := fmt.Sprintf("http://localhost:1114/create-room?num-players=%d", numPlayers)
	got := HttpRequest(t, url, http.StatusOK)

	roomId := got["id"].(string)
	url = fmt.Sprintf("ws://localhost:1114/ws/join?room-id=%s", roomId)
	WsRequest(t, url)
	WsRequest(t, url)
	res := WsRequest(t, url)
	if res.Type != pkg.ServerErrPayloadType || res.Data != "room is full" {
		t.Errorf("got %v, wanted %v", res, pkg.NewErrServerPayload("room is full"))
	}
}
