package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// NOTE: dev only
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{
		hub: hub,
	}
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.hub.Log.Error("Could not upgrade conn", zap.Error(err))
		return
	}

	client := &ws.Client{Conn: conn}
	h.hub.Register <- client

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.hub.Unregister <- client
			h.hub.Log.Info("Client disconnected")
			break
		}
	}
}
