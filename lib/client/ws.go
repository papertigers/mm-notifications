package client

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func (c *Client) StartWatcher(ctx context.Context) <-chan WebSocketEvent {
	out := make(chan WebSocketEvent)

	var headers http.Header
	headers = make(map[string][]string)

	u := url.URL{Scheme: "wss", Host: c.baseurl, Path: "/api/v4/websocket"}
	headers.Set("Cookie", c.cookie)

	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
	}

	// Process messages
	go func() {
		defer ws.Close()
		defer close(out)

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Printf("read: %s\n", err)
				return
			}
			b := bytes.NewBuffer(message)
			var event WebSocketEvent
			err = json.NewDecoder(b).Decode(&event)
			if err != nil {
				log.Printf("decode: %s\n", err)
				return
			}
			out <- event
		}
	}()
	go func() {
		<-ctx.Done()
		err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
		log.Println("Closing websocket")
	}()
	return out
}
