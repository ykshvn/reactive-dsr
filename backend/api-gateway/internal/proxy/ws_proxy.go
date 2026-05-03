package proxy

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// NOTE: dev only
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSProxy struct {
	log *zap.Logger

	// TODO: conn to simulation service
}

func NewWSProxy(l *zap.Logger) *WSProxy {
	return &WSProxy{
		log: l,
	}
}

func (p *WSProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.log.Error("WebSocket upgrade fail", zap.Error(err))
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	p.log.Info("New WS client connected",
		zap.String("Remote addr", r.RemoteAddr),
	)

	// TODO: implement ws logic here
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			p.log.Info("Client disconnected", zap.String("client", r.RemoteAddr))
			break
		}

		p.log.Info("New message from client",
			zap.String("client", r.RemoteAddr),
			zap.String("message", string(message)),
		)

		if err := conn.WriteMessage(messageType, message); err != nil {
			p.log.Error("Could not reply to message")
			break
		}
	}
}
