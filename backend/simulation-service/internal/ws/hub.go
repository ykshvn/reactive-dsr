// Package websocket is required for handling ws conns
package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ykshvn/reactive-dsr/shared/events"
	"go.uber.org/zap"
)

type Client struct {
	Conn *websocket.Conn
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Log        *zap.Logger
	mu         sync.RWMutex
}

func NewHub(l *zap.Logger) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte, 512),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Log:        l,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			h.Log.Info("Client connected to simulation")
		case client := <-h.Unregister:
			h.mu.Lock()
			h.Clients[client] = false
			h.mu.Unlock()
			h.Log.Info("Client disconnected form simulation")
		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.Clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					h.Log.Warn("Could not send message to client", zap.Error(err))
					h.Unregister <- client
				}

			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastEvent(event events.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		h.Log.Error("Could not marshal event", zap.Error(err))
		return
	}
	h.Broadcast <- data
}
