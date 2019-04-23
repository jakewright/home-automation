package routes

import (
	"encoding/json"
	"home-automation/libraries/go/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	//ReadBufferSize:  1024,
	//WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func (c *Controller) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to create websocket upgrader: %v", err)
		return
	}
	defer ws.Close()

	for event := range c.Repository.Events {
		formattedEvent := event.Format()
		b, err := json.Marshal(formattedEvent)
		if err != nil {
			slog.Error("Failed to marshal event: %v", err)
			continue
		}

		if err := ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
			slog.Error("Failed to set write deadline: %v", err)
			return
		}

		if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
			slog.Error("Failed to write message to websocket: %v", err)
			return
		}
	}
}
