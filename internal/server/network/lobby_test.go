package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func request(t *testing.T, url string, expectedStatusCode int) map[string]string {
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
	var got map[string]string
	err = json.Unmarshal(body, &got)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestLobby_CreateRoom(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1111")

	numPlayers := 3
	url := fmt.Sprintf("http://localhost:1111/create-room?num-players=%d", numPlayers)
	got := request(t, url, http.StatusOK)

	if len(got["id"]) != NumIdDigits {
		t.Errorf("got id %s, wanted %d digits", got["id"], NumIdDigits)
	}
}

func TestLobby_CreateRoom_InvalidNumPlayers(t *testing.T) {
	lobby := NewLobby()
	go lobby.Start("1112")

	numPlayers := 0
	url := fmt.Sprintf("http://localhost:1111/create-room?num-players=%d", numPlayers)
	got := request(t, url, http.StatusBadRequest)
	_, ok := got["error"]
	if !ok {
		t.Errorf("expected error response, got %v", got)
	}
}
