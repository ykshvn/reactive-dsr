// Package node describes single node logic
package node

import (
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/domain"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

type NodeActor struct {
	ID           int
	Neighbors    []int
	RouteCache   map[int][]int
	SeenRequests map[string]bool
	Inbox        chan domain.Message
	Hub          *ws.Hub
	log          *zap.Logger
	stop         chan struct{}
}

func NewNodeActor(id int, neighbors []int, hub *ws.Hub, l *zap.Logger) *NodeActor {
	return &NodeActor{
		ID:           id,
		Neighbors:    neighbors,
		RouteCache:   make(map[int][]int),
		SeenRequests: make(map[string]bool),
		Inbox:        make(chan domain.Message, 50),
		Hub:          hub,
		log:          l,
		stop:         make(chan struct{}),
	}
}

func (a *NodeActor) Start() {
	go a.run()
}

func (a *NodeActor) Stop() {
	close(a.stop)
}

func (a *NodeActor) run() {
	for {
		select {
		case msg := <-a.Inbox:
			a.processMessage(msg)
		case <-a.stop:
			return
		}
	}
}

func (a *NodeActor) processMessage(msg domain.Message) {
	switch msg.Type {
	case domain.MessageRREQ:
		a.handleRREQ(msg.RREQ)
	case domain.MessageRREP:
		a.handleRREP(msg.RREP)
	}
}

func (a *NodeActor) handleRREQ(rreq *domain.RREQ) {
	a.log.Info(
		"got RREQ",
		zap.Int("node ID", a.ID),
		zap.Int("from", rreq.Source),
		zap.Int("to", rreq.Destination),
	)
}

func (a *NodeActor) handleRREP(rrep *domain.RREP) {
	a.log.Info(
		"got RREP",
		zap.Int("node ID", a.ID),
		zap.Any("route", rrep.Route),
	)
}
