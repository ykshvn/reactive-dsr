// Package node describes single node logic
package node

import (
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/domain"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

type Node struct {
	ID           int
	Neighbors    []int
	RouteCache   map[int][]int
	SeenRequests map[string]bool
	Inbox        chan domain.Message
	Hub          *ws.Hub
	log          *zap.Logger
	stop         chan struct{}
}

func NewNode(id int, neighbors []int, hub *ws.Hub, l *zap.Logger) *Node {
	return &Node{
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

func (n *Node) Start() {
	go n.run()
}

func (n *Node) Stop() {
	close(n.stop)
}

func (n *Node) run() {
	for {
		select {
		case msg := <-n.Inbox:
			n.processMessage(msg)
		case <-n.stop:
			return
		}
	}
}

func (n *Node) processMessage(msg domain.Message) {
	switch msg.Type {
	case domain.MessageRREQ:
		n.handleRREQ(msg.RREQ)
	case domain.MessageRREP:
		n.handleRREP(msg.RREP)
	}
}

func (n *Node) handleRREQ(rreq *domain.RREQ) {
	n.log.Info(
		"got RREQ",
		zap.Int("node ID", n.ID),
		zap.Int("from", rreq.Source),
		zap.Int("to", rreq.Destination),
	)
}

func (n *Node) handleRREP(rrep *domain.RREP) {
	n.log.Info(
		"got RREP",
		zap.Int("node ID", n.ID),
		zap.Any("route", rrep.Route),
	)
}
