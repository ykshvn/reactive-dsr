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
	simulationServiceURL string
	log                  *zap.Logger
}

func NewWSProxy(simulationServiceURL string, l *zap.Logger) *WSProxy {
	return &WSProxy{
		simulationServiceURL: simulationServiceURL,
		log:                  l,
	}
}

func (p *WSProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.log.Error("WebSocket upgrade fail", zap.Error(err))
	}
	defer clientConn.Close()

	p.log.Info(
		"New WS client connected",
		zap.String("Remote addr", r.RemoteAddr),
	)

	simConn, _, err := websocket.DefaultDialer.Dial(p.simulationServiceURL, nil)
	if err != nil {
		p.log.Error("Failed to connect to simulation service", zap.Error(err))
		return
	}
	defer simConn.Close()

	p.log.Info(
		"WebSocket proxy connection established",
		zap.String("client", r.RemoteAddr),
		zap.String("target", p.simulationServiceURL),
	)

	done := make(chan struct{}, 2)
	go func() {
		copyMessages(clientConn, simConn, "client → simulation")
		done <- struct{}{}
	}()

	go func() {
		copyMessages(simConn, clientConn, "simulation → client")
		done <- struct{}{}
	}()

	<-done
}

func copyMessages(src, dst *websocket.Conn, direction string) {
	for {
		msgType, msg, err := src.ReadMessage()
		if err != nil {
			return
		}

		if err := dst.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
