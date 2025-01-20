package network

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sokmontrey/TicTacToeTuiOnline/pkg"
	"io"
	"net/http"
	"testing"
)

func HttpRequest(t *testing.T, url string, expectedStatusCode int) map[string]any {
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

func WsRequest(t *testing.T, url string) pkg.Payload {
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

func WsRequestWithPayload(t *testing.T, url string, payload pkg.Payload) pkg.Payload {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for {
		conn.WriteJSON(payload)
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
