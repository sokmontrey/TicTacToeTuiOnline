package network

import (
	"fmt"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"net/http"
	"testing"
)

func TestRoom_HandleClient(t *testing.T) {
	lobby := NewLobby()
	port := "2221"
	go lobby.Start(port)

	numPlayers := 2
	url := fmt.Sprintf("http://localhost:%s/create-room?num-players=%d", port, numPlayers)
	got := HttpRequest(t, url, http.StatusOK)

	roomId := got["id"].(string)
	url = fmt.Sprintf("ws://localhost:%s/ws/join?room-id=%s", port, roomId)
	pl := pkg.NewKeypressClientPayload(pkg.KeyCodeConfirm)
	res := WsRequestWithPayload(t, url, pl)
	if res.Type != pkg.ServerOkPayloadType {
		t.Errorf("got %v, wanted %v", res, pkg.NewOkServerPayload("joined room "+roomId))
	}
}
